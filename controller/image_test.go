package controller

import (
	"bytes"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cecobask/hipeople-coding-challenge/util"
)

func init() {
	log.SetOutput(io.Discard)
}

// Removes all images from the uploads folder
func cleanup() {
	files, err := filepath.Glob("../uploads/image-*.*")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func imageUpload(t *testing.T, ctrl ImageController, formKey string) (string, *util.RequestError) {
	// Prepare the request body
	body := new(bytes.Buffer)
	mpWriter := multipart.NewWriter(body)
	part, err := mpWriter.CreateFormFile(formKey, "test.png")
	if err != nil {
		t.Fatal(err)
	}
	fileContents, err := os.ReadFile("../fixtures/test.png")
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
	return ctrl.Upload(req)
}

func TestList(t *testing.T) {
	type test struct {
		name        string
		fixturePath string
		uploadImage bool
	}

	tests := []test{
		{
			name:        "success listing all images",
			fixturePath: "../fixtures/test.png",
			uploadImage: true,
		},
		{
			name:        "empty response when there are no images",
			fixturePath: "",
			uploadImage: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(cleanup)
			ctrl := New()
			// Conditional setup
			if tc.uploadImage == true {
				uploadedImageID, _ := imageUpload(t, ctrl, FORM_KEY_VALUE)
				imagesStr, err := ctrl.List()

				if strings.Contains(imagesStr, uploadedImageID) == false {
					t.Fatalf("Expected result to include the newly uploaded image ID: %v, got: %v", uploadedImageID, imagesStr)
				}
				if err != nil {
					t.Fatalf("Expected error to be nil, got: %v", err.Message)
				}
			} else {
				imagesStr, _ := ctrl.List()
				if imagesStr != "" {
					t.Fatalf("Expected empty response, got: %v", imagesStr)
				}
			}
		})
	}
}

func TestUpload(t *testing.T) {
	type test struct {
		name    string
		formKey string
		err     *util.RequestError
	}

	tests := []test{
		{
			name:    "success uploading an image",
			formKey: FORM_KEY_VALUE,
			err:     nil,
		},
		{
			name:    "error when invalid form key is passed",
			formKey: "invalidKey",
			err: &util.RequestError{
				Status:  http.StatusBadRequest,
				Message: "No form key with value `imageFile` found in the request body",
				Err:     errors.New("invalid form key"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(cleanup)
			ctrl := New()
			fileID, reqErr := imageUpload(t, ctrl, tc.formKey)

			if tc.err != nil {
				if reqErr.Message != tc.err.Message {
					t.Fatalf("Expected error: %v, got: %v", tc.err.Message, reqErr.Message)
				}
			} else {
				if fileID == "" {
					t.Fatal("Expected file ID to not be empty")
				}
			}
		})
	}
}
