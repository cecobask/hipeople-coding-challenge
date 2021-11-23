package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.Handle("/images/", &imageHandler{})
	http.ListenAndServe(":8080", mux)
}
