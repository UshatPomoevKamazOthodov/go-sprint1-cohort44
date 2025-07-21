package handlers

import (
	"encoding/json"
	storage "go-sprint1-cohort44/internal/cache"
	"net/http"
	"net/url"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func PostJSONHandle(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Missing 'url' field", http.StatusBadRequest)
		return
	}

	urlPair := storage.GlobalStorage.Save(req.URL)

	base := "http://" + r.Host
	fullURL, err := url.JoinPath(base, urlPair.URLReduced)
	if err != nil {
		http.Error(w, "Failed to build URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := Response{
		Result: fullURL,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to write JSON", http.StatusInternalServerError)
	}
}
