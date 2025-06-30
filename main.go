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

	config := cfg.GetConfigData()

	// Выводим информацию о конфигурации
	fmt.Println("Server on: " + config.ServerAddr)
	fmt.Println("Base URL: " + config.BaseURL)
	// Создаем новый маршрутизатор
	r := mux.NewRouter()

	// Регистрируем обработчики
	r.HandleFunc("/getUrl", middleware.GetUrlHandle).Methods("GET")
	r.HandleFunc("/postUrl", middleware.PostUrlHandle).Methods("POST")

	// Запускаем сервер
	err := http.ListenAndServe(config.ServerAddr, r)
	if err != nil {
		panic(err)
	}
}
