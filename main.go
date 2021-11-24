package main

import (
	"log"
	"net/http"

	"github.com/cecobask/hipeople-coding-challenge/controller"
	"github.com/cecobask/hipeople-coding-challenge/handler"
)

func main() {
	controller := controller.New()
	handler := handler.New(controller)
	mux := http.NewServeMux()
	mux.Handle("/images/", handler)
	log.Fatal(http.ListenAndServe(":8080", mux))
	log.Println("Go server started at port 8080")
}
