package main

import (
	"github.com/go-chi/chi/v5"
	"go-sprint1-cohort44/internal/cache"
	"go-sprint1-cohort44/internal/cfg"
	"go-sprint1-cohort44/internal/handlers"
	"go-sprint1-cohort44/internal/middleware"
	"log"
	"net/http"
)

func main() {
	config := cfg.GetConfigData()
	storage.InitGlobalStorage()

	// Выводим информацию о конфигурации
	log.Printf("Server on: " + config.ServerAddr)
	log.Printf("Base URL: " + config.BaseURL)

	r := chi.NewRouter()

	r.Use(middleware.Logger())

	// Регистрируем обработчики
	r.Get("/getUrl", handlers.GetUrlHandle)
	r.Post("/postUrl", handlers.PostUrlHandle)

	// Запускаем сервер
	err := http.ListenAndServe(config.ServerAddr, r)
	if err != nil {
		log.Fatal(err)
	}
}
