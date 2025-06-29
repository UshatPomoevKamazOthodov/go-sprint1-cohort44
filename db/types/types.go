package types

import (
	"gopkg.in/mgo.v2/bson"
)

type UrlAddress struct {
	Id         bson.ObjectId `bson:"_id"`
	Url        string        `bson:"url"`
	UrlReduced string        `bson:"url_reduced"`
}
