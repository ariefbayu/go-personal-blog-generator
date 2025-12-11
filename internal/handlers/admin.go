package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func ServeDashboard(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "personal-blog-generator", "admin-files", "index.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Admin dashboard template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}