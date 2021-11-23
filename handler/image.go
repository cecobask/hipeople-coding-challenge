package handler

import (
	"log"
	"net/http"
	"regexp"

	"github.com/cecobask/hipeople-coding-challenge/controller"
)

var (
	listImagesRegex  = regexp.MustCompile(`^\/images[\/]*$`)
	getImageRegex    = regexp.MustCompile(`^\/images\/(\d+)$`)
	uploadImageRegex = regexp.MustCompile(`^\/images[\/]*$`)
)

type ImageHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	Upload(http.ResponseWriter, *http.Request)
}

type imageHandler struct {
	imageController controller.ImageController
}

func New(ctrl controller.ImageController) ImageHandler {
	return &imageHandler{
		imageController: ctrl,
	}
}

func (h *imageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && listImagesRegex.MatchString(r.URL.Path):
		h.List(w, r)
		return
	case r.Method == http.MethodGet && getImageRegex.MatchString(r.URL.Path):
		h.GetByID(w, r)
		return
	case r.Method == http.MethodPost && uploadImageRegex.MatchString(r.URL.Path):
		h.Upload(w, r)
		return
	default:
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}
}

func (h *imageHandler) List(w http.ResponseWriter, r *http.Request) {
	log.Println("List images called")
	fileIDs, err := h.imageController.List()
	if err != nil {
		http.Error(w, err.Message, err.Status)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fileIDs))
}

func (h *imageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Get image by ID called")
	w.WriteHeader(http.StatusOK)
}

func (h *imageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload image called")
	fileID, err := h.imageController.Upload(r)
	if err != nil {
		http.Error(w, err.Message, err.Status)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fileID))
}
