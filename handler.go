package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const FORM_KEY_VALUE = "imageFile"

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
	// Match pattern for uploaded image files
	files, err := filepath.Glob("uploads/image-*.*")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Number of images found", len(files))

	// Only return the unique idenfitier for each file, which can be used for retrieving by ID
	var fileIDs []string
	for _, file := range files {
		startIndex := strings.IndexByte(file, '-') + 1
		endIndex := strings.IndexByte(file, '.')
		fileID := file[startIndex:endIndex]
		fileIDs = append(fileIDs, fileID)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.Join(fileIDs, ",")))
}

func (h *imageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Get image by ID called")
	w.WriteHeader(http.StatusOK)
}

func (h *imageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	log.Println("Upload image called")
	// Max file size 10 MB
	r.ParseMultipartForm(10 << 20)
	// Retrieve the first file for the `imageFile` form key
	file, _, err := r.FormFile(FORM_KEY_VALUE)
	if err != nil {
		log.Println(err)
		http.Error(w, "No form key with value `imageFile` found in the request body", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary file within the `uploads` directory that follows a naming pattern
	tempFile, err := os.CreateTemp("uploads", "image-*.png")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Read the contents of the uploaded image into a byte slice and write it to the temporary file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tempFile.Write(fileBytes)

	// Return the unique ID of the image when successful
	filePath := tempFile.Name()
	startIndex := strings.IndexByte(filePath, '-') + 1
	endIndex := strings.IndexByte(filePath, '.')
	uniqueID := filePath[startIndex:endIndex]
	w.WriteHeader(http.StatusOK)
	log.Printf("Successfully uploaded file with ID %s", uniqueID)
	w.Write([]byte(uniqueID))
}
