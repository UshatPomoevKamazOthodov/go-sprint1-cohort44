package db

import (
	"context"
	"errors"
	"fmt"
	"go-sprint1-cohort44/cfg"
	"go-sprint1-cohort44/db/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func ConnectToMongo() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("url_reduce").Collection("url_reducer")
	return collection
}

func GetUrl(urlToFind string) string {
	collection := ConnectToMongo()
	var result types.UrlAddress

	query := bson.M{
		"url_reduced": bson.M{
			"$eq": urlToFind,
		},
	}

	err := collection.FindOne(context.TODO(), query).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return "not_found"
	}
	if err != nil {
		log.Fatal(err)
	}
	return result.Url
}

func InsertUrl(urlToEncode string) (string, error) {
	checkUrl := GetUrl(urlToEncode)
	if checkUrl != "not_found" {
		return "", errors.New("url already exists")
	}

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	urlEncoded := string(b)
	objId := bson.NewObjectId()
	cfg := cfg.GetConfigData()

	insert := types.UrlAddress{
		Id:         objId,
		Url:        cfg.BaseURL + urlToEncode,
		UrlReduced: urlEncoded,
	}

	collection := ConnectToMongo()
	insertResult, err := collection.InsertOne(context.TODO(), insert)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return urlEncoded, nil
}
