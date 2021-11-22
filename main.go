package main

import (
	"net/http"
	"regexp"
)

var (
	listImagesRegex  = regexp.MustCompile(`^\/images[\/]*$`)
	getImageRegex    = regexp.MustCompile(`^\/images\/(\d+)$`)
	uploadImageRegex = regexp.MustCompile(`^\/images[\/]*$`)
)

type imageHandler struct{}

func (h *imageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listImagesRegex.MatchString(r.URL.Path):
		h.GetAll(w, r)
		return
	case r.Method == http.MethodGet && getImageRegex.MatchString(r.URL.Path):
		h.GetByID(w, r)
		return
	case r.Method == http.MethodPost && uploadImageRegex.MatchString(r.URL.Path):
		h.Upload(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

func (h *imageHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("get all images called"))
}

func (h *imageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("get image by id called"))
}

func (h *imageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("upload image called"))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("route not found"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/images/", &imageHandler{})
	http.ListenAndServe(":8080", mux)
}
