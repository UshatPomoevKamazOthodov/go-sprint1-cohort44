package middleware

import (
	"go-sprint1-cohort44/cfg"
	"go-sprint1-cohort44/db"
	"net/http"
)

func PostUrlHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	config := cfg.GetConfigData()
	shortenedUrl, err := db.InsertUrl(r.FormValue("url"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(201)
	w.Header().Set("Location", config.BaseURL+shortenedUrl)
	return
}
