package handler

import (
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/cecobask/hipeople-coding-challenge/controller"
)

var (
	listImagesRegex  = regexp.MustCompile(`^\/images[\/]*$`)
	getImageRegex    = regexp.MustCompile(`^\/images\/(\d+)$`)
	uploadImageRegex = regexp.MustCompile(`^\/images[\/]*$`)
)

// ImageHandler methods
type ImageHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	Upload(http.ResponseWriter, *http.Request)
}

type imageHandler struct {
	imageController controller.ImageController
}

// New ImageHandler constructor
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
	imageIDs, err := h.imageController.List()
	if err != nil {
		http.Error(w, err.Message, err.Status)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(imageIDs))
}

func (h *imageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Get image by ID called")
	imageID := path.Base(r.URL.Path)
	imageBytes, err := h.imageController.GetByID(imageID)
	if err != nil {
		http.Error(w, err.Message, err.Status)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(imageBytes)
}

func (h *imageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload image called")
	imageID, err := h.imageController.Upload(r)
	if err != nil {
		http.Error(w, err.Message, err.Status)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(imageID))
}
