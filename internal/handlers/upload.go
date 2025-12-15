package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	maxUploadSize = 5 << 20 // 5MB
)

var (
	uploadDir = ""
)

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// UploadImageHandler handles image file uploads for the WYSIWYG editor
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	outputPath := os.Getenv("OUTPUT_PATH")
	if outputPath == "" {
		http.Error(w, "OUTPUT_PATH not configured", http.StatusInternalServerError)
		return
	}
	uploadDir = filepath.Join(outputPath, "images")

	// Parse multipart form with max memory
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	if header.Size > maxUploadSize {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Validate content type
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	file.Seek(0, 0) // Reset file pointer
	contentType := http.DetectContentType(buffer[:n])
	if !allowedTypes[contentType] {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		http.Error(w, "Invalid file extension", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	id := uuid.New().String()
	filename := id + ext

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	// Create destination file
	dstPath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Return JSON response with image URL
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"data": {"filePath": "/images/%s"}}`, filename)
}
