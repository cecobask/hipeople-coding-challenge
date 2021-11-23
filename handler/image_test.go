package handler

import (
	"errors"
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
	type test struct {
		name         string
		imageID      string
		expectedCode int
		err          *util.RequestError
	}

	tests := []test{
		{
			name:         "success listing all images",
			imageID:      "1637696279",
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
			req, res := setupHttpRequest(t, http.MethodGet, "/images/", nil)
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

func setupHttpRequest(t *testing.T, method, path string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	return req, res
}

func TestUpload(t *testing.T) {
	type test struct {
		name         string
		imageID      string
		expectedCode int
		err          *util.RequestError
	}

	tests := []test{
		{
			name:         "success uploading an image",
			imageID:      "1637696279",
			expectedCode: http.StatusOK,
			err:          nil,
		},
		{
			name:         "error when invalid form key is passed",
			imageID:      "1637696279",
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

			req, res := setupHttpRequest(t, http.MethodPost, "/images/", nil)
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
