package handlers

import (
	"go-sprint1-cohort44/internal/cache"
	"go-sprint1-cohort44/internal/cfg"
	"net/http"
)

func GetUrlHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	config := cfg.GetConfigData()
	queryParams := r.URL.Query()
	urlParam := queryParams.Get("url")

	original, found := storage.GlobalStorage.GetUrl(urlParam)
	if !found {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + config.ServerAddr + "/" + original))
	return
}
