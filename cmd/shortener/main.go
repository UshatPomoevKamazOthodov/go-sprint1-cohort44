package main

import (
	"github.com/go-chi/chi/v5"
	"go-sprint1-cohort44/internal/cache"
	"go-sprint1-cohort44/internal/cfg"
	"go-sprint1-cohort44/internal/middleware"
	"log"
	"net/http"
)

func main() {
	config := cfg.GetConfigData()
	cache.InitCache()

	// Выводим информацию о конфигурации
	log.Printf("Server on: " + config.ServerAddr)
	log.Printf("Base URL: " + config.BaseURL)

	r := chi.NewRouter()

	// Регистрируем обработчики
	r.Get("/getUrl", middleware.GetUrlHandle)
	r.Post("/postUrl", middleware.PostUrlHandle)

	// Запускаем сервер
	err := http.ListenAndServe(config.ServerAddr, r)
	if err != nil {
		panic(err)
	}
}
