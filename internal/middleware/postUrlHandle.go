package middleware

import (
	"fmt"
	"go-sprint1-cohort44/internal/cache"
	"go-sprint1-cohort44/internal/cfg"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func PostUrlHandle(w http.ResponseWriter, r *http.Request) {
	config := cfg.GetConfigData()
	if config == nil {
		http.Error(w, "Configuration is not initialized", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	postData := r.FormValue("url")
	if postData == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	// Проверяем кэш
	checkUrl, found := cache.GlobalCache.Get(postData)
	log.Println(cache.GlobalCache.Get(postData))
	if found {
		if str, ok := checkUrl.(string); ok {
			http.Error(w, fmt.Sprintf("%s already exists", str), http.StatusBadRequest)
			return
		} else {
			http.Error(w, "Invalid value in cache", http.StatusInternalServerError)
			return
		}
	}

	shortenedURL := randomString(10)
	cache.GlobalCache.Set(postData, shortenedURL, 24*time.Hour)
	log.Println(cache.GlobalCache.Get(postData))

	location := "http://" + config.ServerAddr + "/" + shortenedURL

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("Shortened URL: " + shortenedURL))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
