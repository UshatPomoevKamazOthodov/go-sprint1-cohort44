package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"go-sprint1-cohort44/cfg"
	"go-sprint1-cohort44/middleware"
	"net/http"
)

func main() {
	flag.Parse()

	cfg := cfg.GetConfigData()

	// Выводим информацию о конфигурации
	fmt.Println("Server on: " + cfg.ServerAddr)
	fmt.Println("Base URL: " + cfg.BaseURL)
	// Создаем новый маршрутизатор
	r := mux.NewRouter()

	// Регистрируем обработчики
	r.HandleFunc("/getUrl", middleware.GetUrlHandle).Methods("GET")
	r.HandleFunc("/postUrl", middleware.PostUrlHandle).Methods("POST")

	// Запускаем сервер
	err := http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		panic(err)
	}
}
