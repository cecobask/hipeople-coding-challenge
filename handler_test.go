package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func init() {
	log.SetOutput(io.Discard)
}

// Removes all images from the uploads folder
func cleanup() {
	files, err := filepath.Glob("uploads/image-*.*")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func imageUpload(t *testing.T, formKey string) *httptest.ResponseRecorder {
	// Prepare the request body
	body := new(bytes.Buffer)
	mpWriter := multipart.NewWriter(body)
	part, err := mpWriter.CreateFormFile(formKey, "test.png")
	if err != nil {
		t.Fatal(err)
	}
	fileContents, err := os.ReadFile("fixtures/test.png")
	if err != nil {
		t.Fatal(err)
	}
	part.Write(fileContents)
	mpWriter.Close()

	// Create the request/response structs and call the image upload handler
	req, err := http.NewRequest(http.MethodPost, "/images/", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	resRecorder := httptest.NewRecorder()
	handler := &imageHandler{}
	handler.Upload(resRecorder, req)
	return resRecorder
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
			formKey:      FORM_KEY_VALUE,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(cleanup)
			resRecorder := imageUpload(t, tc.formKey)

			// Passing criteria
			if status := resRecorder.Code; status != tc.expectedCode {
				t.Fatalf("Expected status code: %v, got: %v", tc.expectedCode, status)
			}
		})
	}
}

func setupListRequest(t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
	// Create the request/response structs and call the list images handler
	req, err := http.NewRequest(http.MethodGet, "/images/", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	return req, res
}

func TestList(t *testing.T) {
	type test struct {
		name         string
		fixturePath  string
		expectedCode int
		doSetup      bool
	}

	tests := []test{
		{
			name:         "success listing all images",
			fixturePath:  "fixtures/test.png",
			expectedCode: http.StatusOK,
			doSetup:      true,
		},
		{
			name:         "empty response when there are no images",
			fixturePath:  "",
			expectedCode: http.StatusOK,
			doSetup:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(cleanup)
			handler := &imageHandler{}
			// Conditional setup
			if tc.doSetup == true {
				resUpload := imageUpload(t, FORM_KEY_VALUE)
				uploadedImageID := resUpload.Body.String()
				req, resList := setupListRequest(t)
				handler.List(resList, req)
				imagesStr := resList.Body.String()

				// Passing criteria
				if status := resList.Code; status != tc.expectedCode {
					t.Fatalf("Expected status code: %v, got: %v", tc.expectedCode, status)
				}
				if strings.Contains(imagesStr, uploadedImageID) == false {
					t.Fatalf("Expected result to include the newly uploaded image ID: %v, got: %v", uploadedImageID, imagesStr)
				}
			} else {
				req, resList := setupListRequest(t)
				handler.List(resList, req)
				imagesStr := resList.Body.String()

				// Passing criteria
				if status := resList.Code; status != tc.expectedCode {
					t.Fatalf("Expected status code: %v, got: %v", tc.expectedCode, status)
				}
				if imagesStr != "" {
					t.Fatalf("Expected empty response, got: %v", imagesStr)
				}
			}
		})
	}
}
