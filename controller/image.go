package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cecobask/hipeople-coding-challenge/util"
)

const FORM_KEY_VALUE = "imageFile"

//go:generate mockgen -destination=mock/image.go -package=mock github.com/cecobask/hipeople-coding-challenge/controller ImageController
type ImageController interface {
	List() (string, *util.RequestError)
	GetByID(imageID string) ([]byte, *util.RequestError)
	Upload(r *http.Request) (string, *util.RequestError)
}

type imageController struct{}

func New() ImageController {
	return &imageController{}
}

func (c *imageController) List() (string, *util.RequestError) {
	// Match pattern for uploaded image files
	files, err := filepath.Glob("uploads/image-*.*")
	if err != nil {
		log.Println(err)
		return "", &util.RequestError{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
			Err:     err,
		}
	}
	log.Println("Number of images found", len(files))

	// Only return the unique idenfitier for each file, which can be used for retrieving by ID
	var fileIDs []string
	for _, file := range files {
		startIndex := strings.IndexByte(file, '-') + 1
		endIndex := strings.LastIndex(file, ".")
		fileID := file[startIndex:endIndex]
		fileIDs = append(fileIDs, fileID)
	}
	return strings.Join(fileIDs, ","), nil
}

func (c *imageController) GetByID(imageID string) ([]byte, *util.RequestError) {
	// Look for the specified image file
	images, err := filepath.Glob(fmt.Sprintf("uploads/image-%s.*", imageID))
	if err != nil {
		log.Println(err)
		return nil, &util.RequestError{
			Status:  http.StatusBadRequest,
			Message: "Pattern malformed",
			Err:     err,
		}
	}
	if images == nil {
		log.Println(err)
		return nil, &util.RequestError{
			Status:  http.StatusNotFound,
			Message: "Image not found",
			Err:     err,
		}
	}

	// Return the specified file
	imageBytes, err := os.ReadFile(images[0])
	if err != nil {
		log.Println(err)
		return nil, &util.RequestError{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
			Err:     err,
		}
	}
	log.Println("Successfully retrieved image with ID", imageID)
	return imageBytes, nil
}

func (c *imageController) Upload(r *http.Request) (string, *util.RequestError) {
	// Max file size 10 MB
	r.ParseMultipartForm(10 << 20)
	// Retrieve the first file for the `imageFile` form key
	file, _, err := r.FormFile(FORM_KEY_VALUE)
	if err != nil {
		log.Println(err)
		return "", &util.RequestError{
			Status:  http.StatusBadRequest,
			Message: "No form key with value `imageFile` found in the request body",
			Err:     err,
		}
	}
	defer file.Close()

	// Create a temporary file within the `uploads` directory that follows a naming pattern
	tempFile, err := os.CreateTemp("uploads", "image-*.png")
	if err != nil {
		log.Println(err)
		return "", &util.RequestError{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
			Err:     err,
		}
	}
	defer tempFile.Close()

	// Read the contents of the uploaded image into a byte slice and write it to the temporary file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return "", &util.RequestError{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
			Err:     err,
		}
	}
	tempFile.Write(fileBytes)

	// Return the unique ID of the image when successful
	filePath := tempFile.Name()
	startIndex := strings.IndexByte(filePath, '-') + 1
	endIndex := strings.LastIndex(filePath, ".")
	uniqueID := filePath[startIndex:endIndex]
	log.Printf("Successfully uploaded file with ID %s", uniqueID)

	return uniqueID, nil
}
