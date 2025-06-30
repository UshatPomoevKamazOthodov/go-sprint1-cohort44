package middleware

import (
	"go-sprint1-cohort44/cfg"
	"go-sprint1-cohort44/db"
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
	shortenedUrl := db.GetUrl(urlParam)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + config.ServerAddr + shortenedUrl))
	return
}
