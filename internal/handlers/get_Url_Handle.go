package handlers

import (
	"go-sprint1-cohort44/internal/cache"
	"go-sprint1-cohort44/internal/cfg"
	"net/http"
)

func GetUrlHandle(w http.ResponseWriter, r *http.Request) {
	config := cfg.GetConfigData()
	GlobalStorage := storage.InitGlobalStorage(config.BasePathToFile)
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()
	urlParam := queryParams.Get("url")

	original, found := GlobalStorage.GetUrl(urlParam)
	if !found {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + config.ServerAddr + "/" + original))
	return
}
