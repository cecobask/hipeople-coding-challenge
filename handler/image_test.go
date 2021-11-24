package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cecobask/hipeople-coding-challenge/controller/mock"
	"github.com/cecobask/hipeople-coding-challenge/util"
	"github.com/golang/mock/gomock"
)

const imageID = "1637696279"

type testCase struct {
	name         string
	imageID      string
	imageBytes   []byte
	expectedCode int
	err          *util.RequestError
}

func init() {
	log.SetOutput(io.Discard)
}

func setupMocks(t *testing.T) (ImageHandler, *mock.MockImageController) {
	mockCtrl := gomock.NewController(t)
	mockImageCtrl := mock.NewMockImageController(mockCtrl)
	handler := New(mockImageCtrl)
	return handler, mockImageCtrl
}

func TestList(t *testing.T) {
	tests := []testCase{
		{
			name:         "success listing all images",
			imageID:      imageID,
			expectedCode: http.StatusOK,
			err:          nil,
		},
		{
			name:         "error in the controller",
			imageID:      "",
			expectedCode: http.StatusInternalServerError,
			err: &util.RequestError{
				Status:  http.StatusInternalServerError,
				Message: "controller error",
				Err:     errors.New("controller error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler, mockImageCtrl := setupMocks(t)
			mockImageCtrl.EXPECT().
				List().
				Return(tc.imageID, tc.err).
				Times(1)
			req, res := setupHTTPRequest(t, http.MethodGet, "/images/", nil)
			handler.List(res, req)

			if status := res.Code; status != tc.expectedCode {
				t.Fatalf("Expected status code: %v, got: %v", tc.expectedCode, status)
			}
			if tc.err != nil {
				if imagesStr := res.Body.String(); !strings.Contains(imagesStr, tc.err.Message) {
					t.Fatalf("Expected request error: %v, got: %v", tc.err.Message, imagesStr)
				}
			} else {
				if imagesStr := res.Body.String(); strings.Contains(imagesStr, tc.imageID) == false {
					t.Fatalf("Expected result to include the image ID: %s, got: %s", tc.imageID, imagesStr)
				}
			}
		})
	}
}

func setupHTTPRequest(t *testing.T, method, path string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	return req, res
}

func TestUpload(t *testing.T) {
	tests := []testCase{
		{
			name:         "success uploading an image",
			imageID:      imageID,
			expectedCode: http.StatusOK,
			err:          nil,
		},
		{
			name:         "error in the controller",
			imageID:      imageID,
			expectedCode: http.StatusInternalServerError,
			err: &util.RequestError{
				Status:  http.StatusInternalServerError,
				Message: "controller error",
				Err:     errors.New("controller error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler, mockImageCtrl := setupMocks(t)
			mockImageCtrl.EXPECT().
				Upload(gomock.Any()).
				Return(tc.imageID, tc.err).
				Times(1)

			req, res := setupHTTPRequest(t, http.MethodPost, "/images/", nil)
			handler.Upload(res, req)

			if status := res.Code; status != tc.expectedCode {
				t.Fatalf("Expected status code: %v, got: %v", tc.expectedCode, status)
			}
			if tc.err != nil {
				if imagesStr := res.Body.String(); !strings.Contains(imagesStr, tc.err.Message) {
					t.Fatalf("Expected request error: %v, got: %v", tc.err.Message, imagesStr)
				}
			} else {
				if imagesStr := res.Body.String(); strings.Contains(imagesStr, tc.imageID) == false {
					t.Fatalf("Expected result to include the image ID: %s, got: %s", tc.imageID, imagesStr)
				}
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	tests := []testCase{
		{
			name:         "success retrieving an image",
			imageBytes:   []byte{1, 2, 3},
			expectedCode: http.StatusOK,
			err:          nil,
		},
		{
			name:         "error in the controller",
			imageBytes:   nil,
			expectedCode: http.StatusNotFound,
			err: &util.RequestError{
				Status:  http.StatusNotFound,
				Message: "controller error",
				Err:     errors.New("controller error"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler, mockImageCtrl := setupMocks(t)
			mockImageCtrl.EXPECT().
				GetByID(gomock.Any()).
				Return(tc.imageBytes, tc.err).
				Times(1)

			req, res := setupHTTPRequest(t, http.MethodGet, fmt.Sprintf("/images/%s", imageID), nil)
			handler.GetByID(res, req)

			if status := res.Code; status != tc.expectedCode {
				t.Fatalf("Expected status code: %v, got: %v", tc.expectedCode, status)
			}
			if tc.err != nil {
				if imagesStr := res.Body.String(); !strings.Contains(imagesStr, tc.err.Message) {
					t.Fatalf("Expected request error: %v, got: %v", tc.err.Message, imagesStr)
				}
			}
		})
	}
}
