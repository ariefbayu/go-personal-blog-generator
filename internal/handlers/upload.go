package handlers

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
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

	if OutputPath == "" {
		http.Error(w, "OUTPUT_PATH not configured", http.StatusInternalServerError)
		return
	}

	slug := r.URL.Query().Get("slug")
	if slug != "" {
		uploadDir = filepath.Join(OutputPath, "images", slug)
	} else {
		uploadDir = filepath.Join(OutputPath, "images")
	}

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

	// Reset file pointer to generate thumbnail
	file.Seek(0, 0)

	// Decode source image
	srcImg, format, err := image.Decode(file)
	if err == nil {
		// Generate thumbnail
		thumbImg := createThumbnail(srcImg, 300)

		// Create thumbnail file
		thumbFilename := id + "_thumb" + ext
		if ext == ".webp" {
			thumbFilename = id + "_thumb.jpg"
		}
		thumbPath := filepath.Join(uploadDir, thumbFilename)
		thumbFile, err := os.Create(thumbPath)
		if err == nil {
			defer thumbFile.Close()
			// Encode based on format
			switch format {
			case "png":
				png.Encode(thumbFile, thumbImg)
			case "gif":
				gif.Encode(thumbFile, thumbImg, nil)
			default: // jpeg, webp, etc.
				jpeg.Encode(thumbFile, thumbImg, &jpeg.Options{Quality: 85})
			}
		}
	}

	// Return JSON response with image URL and thumbnail URL
	w.Header().Set("Content-Type", "application/json")
	var filenamePath, thumbnailPath string
	if slug != "" {
		filenamePath = fmt.Sprintf("/images/%s/%s", slug, filename)
		if ext == ".webp" {
			thumbnailPath = fmt.Sprintf("/images/%s/%s_thumb.jpg", slug, id)
		} else {
			thumbnailPath = fmt.Sprintf("/images/%s/%s_thumb%s", slug, id, ext)
		}
	} else {
		filenamePath = fmt.Sprintf("/images/%s", filename)
		if ext == ".webp" {
			thumbnailPath = fmt.Sprintf("/images/%s_thumb.jpg", id)
		} else {
			thumbnailPath = fmt.Sprintf("/images/%s_thumb%s", id, ext)
		}
	}
	fmt.Fprintf(w, `{"data": {"filePath": "%s", "thumbnailPath": "%s"}}`, filenamePath, thumbnailPath)
}

func resizeImage(img image.Image, width, height int) image.Image {
	srcBounds := img.Bounds()
	dstBounds := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(dstBounds)

	dx := float64(srcBounds.Dx()) / float64(width)
	dy := float64(srcBounds.Dy()) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x)*dx) + srcBounds.Min.X
			srcY := int(float64(y)*dy) + srcBounds.Min.Y
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}

func createThumbnail(img image.Image, maxDim int) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	var newW, newH int
	if w > h {
		newW = maxDim
		newH = int(float64(h) * float64(maxDim) / float64(w))
	} else {
		newH = maxDim
		newW = int(float64(w) * float64(maxDim) / float64(h))
	}

	return resizeImage(img, newW, newH)
}
