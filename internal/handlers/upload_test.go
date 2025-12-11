package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestUploadImageHandler(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		contentType    string
		content        []byte
		expectedStatus int
		expectFile     bool
	}{
		{
			name:           "valid JPEG",
			filename:       "test.jpg",
			contentType:    "image/jpeg",
			content:        []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}, // JPEG header
			expectedStatus: http.StatusOK,
			expectFile:     true,
		},
		{
			name:           "valid PNG",
			filename:       "test.png",
			contentType:    "image/png",
			content:        []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG header
			expectedStatus: http.StatusOK,
			expectFile:     true,
		},
		{
			name:           "invalid content type",
			filename:       "test.txt",
			contentType:    "text/plain",
			content:        []byte("text content"),
			expectedStatus: http.StatusBadRequest,
			expectFile:     false,
		},
		{
			name:           "invalid extension",
			filename:       "test.exe",
			contentType:    "image/jpeg",
			content:        []byte("fake jpeg content"),
			expectedStatus: http.StatusBadRequest,
			expectFile:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a unique temporary directory for this test
			tempDir, err := os.MkdirTemp("", "upload_test_"+strings.ReplaceAll(tt.name, " ", "_"))
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			// Override uploadDir for this test
			originalUploadDir := uploadDir
			uploadDir = tempDir
			defer func() { uploadDir = originalUploadDir }()
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("image", tt.filename)
			if err != nil {
				t.Fatal(err)
			}
			part.Write(tt.content)

			// Set content type for the part
			if tt.contentType != "" {
				// Note: In real multipart, this would be set automatically, but for testing we need to simulate
				// Actually, let's check what the header contains
			}
			writer.Close()

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/api/upload/image", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			UploadImageHandler(w, req)

			// Check status
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectFile {
				// Check if file was created
				files, err := os.ReadDir(tempDir)
				if err != nil {
					t.Fatal(err)
				}
				if len(files) != 1 {
					t.Errorf("expected 1 file, got %d", len(files))
				}

				// Check response contains file path
				response := w.Body.String()
				if !strings.Contains(response, "/images/") {
					t.Errorf("response should contain /images/ path, got: %s", response)
				}
			} else {
				// Check no file was created
				files, err := os.ReadDir(tempDir)
				if err != nil {
					t.Fatal(err)
				}
				if len(files) > 0 {
					t.Errorf("expected no files, got %d", len(files))
				}
			}
		})
	}
}

func TestUploadImageHandler_FileTooLarge(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "upload_test_large")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Override uploadDir for testing
	originalUploadDir := uploadDir
	uploadDir = tempDir
	defer func() { uploadDir = originalUploadDir }()

	// Create large content (6MB > 5MB limit)
	largeContent := make([]byte, 6<<20)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", "large.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write(largeContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	UploadImageHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}