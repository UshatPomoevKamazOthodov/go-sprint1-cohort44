package main

import (
	"go-sprint1-cohort44/middleware"
	"net/http"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	http.HandleFunc(`/getUrl`, middleware.GetUrlHandle)
	http.HandleFunc(`/postUrl`, middleware.PostUrlHandle)
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	} else {
		println("Listening on port 3000")
	}
}
