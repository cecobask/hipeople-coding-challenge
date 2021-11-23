package main

import (
	"net/http"

	"github.com/cecobask/hipeople-coding-challenge/controller"
	"github.com/cecobask/hipeople-coding-challenge/handler"
)

func main() {
	controller := controller.New()
	handler := handler.New(controller)
	mux := http.NewServeMux()
	mux.Handle("/images/", handler)
	http.ListenAndServe(":8080", mux)
}
