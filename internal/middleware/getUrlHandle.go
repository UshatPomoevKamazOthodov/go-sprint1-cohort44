package middleware

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

	savedUrl, a := cache.GlobalCache.Get(urlParam)
	if !a {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + config.ServerAddr + "/" + savedUrl.(string)))
	return
}
