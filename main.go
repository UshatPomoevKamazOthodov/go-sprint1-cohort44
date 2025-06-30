package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-sprint1-cohort44/middleware"
	"net/http"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	// Создаем новый маршрутизатор
	r := mux.NewRouter()

	// Регистрируем обработчики
	r.HandleFunc("/getUrl", middleware.GetUrlHandle).Methods("GET")
	r.HandleFunc("/postUrl", middleware.PostUrlHandle).Methods("POST")

	// Запускаем сервер
	fmt.Println("Listening on port :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
