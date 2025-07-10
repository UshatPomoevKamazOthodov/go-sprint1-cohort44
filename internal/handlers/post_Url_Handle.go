package handlers

import (
	"fmt"
	"go-sprint1-cohort44/internal/cache"
	"go-sprint1-cohort44/internal/cfg"
	"log"
	"net/http"
	"net/url"
)

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

	checkForUrl, found := storage.GlobalStorage.GetByOriginal(postData)
	if found {
		http.Error(w, fmt.Sprintf("%s already exists", checkForUrl.URLReduced), http.StatusBadRequest)
		return
	}

	urlPair := storage.GlobalStorage.Save(postData)

	location, err := url.JoinPath("http://"+config.ServerAddr, urlPair.URLReduced)
	if err != nil {
		log.Printf("Server exited with error: %v", err)
	}
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Shortened URL: " + urlPair.URLReduced))
	if err != nil {
		log.Printf("Failed to write response %v", err)
	}
}
