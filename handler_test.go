package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

// Removes all images from the uploads folder
func cleanup() {
	files, err := filepath.Glob("uploads/image-*")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func TestUpload(t *testing.T) {
	type test struct {
		name         string
		formKey      string
		expectedCode int
	}

	tests := []test{
		{
			name:         "error when invalid form key is passed",
			formKey:      "invalidKey",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "success uploading an image",
			formKey:      "imageFile",
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Prepare the request body
			body := new(bytes.Buffer)
			mpWriter := multipart.NewWriter(body)
			part, err := mpWriter.CreateFormFile(test.formKey, "test.png")
			if err != nil {
				t.Fatal(err)
			}
			fileContents, err := ioutil.ReadFile("fixtures/test.png")
			if err != nil {
				t.Fatal(err)
			}
			part.Write(fileContents)
			mpWriter.Close()
			defer cleanup()

			// Create the request/response structs and call the image upload handler
			req, err := http.NewRequest(http.MethodPost, "/images/", body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", mpWriter.FormDataContentType())
			resRecorder := httptest.NewRecorder()
			handler := &imageHandler{}
			handler.Upload(resRecorder, req)

			// Passing criteria
			if status := resRecorder.Code; status != test.expectedCode {
				t.Fatalf("Expected status code: %v, got: %v", test.expectedCode, status)
			}
		})
	}
}
