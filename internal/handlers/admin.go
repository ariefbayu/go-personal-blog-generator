package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var AdminFilesPath string

func init() {
	AdminFilesPath = os.Getenv("ADMIN_FILES_PATH")
	if AdminFilesPath == "" {
		// Default to user directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory
			AdminFilesPath = "./admin-files"
		} else {
			AdminFilesPath = filepath.Join(homeDir, ".personal-blog-generator", "admin-files")
		}
	} else if strings.HasPrefix(AdminFilesPath, "~") {
		// Expand ~ to home directory
		homeDir, err := os.UserHomeDir()
		if err == nil {
			AdminFilesPath = strings.Replace(AdminFilesPath, "~", homeDir, 1)
		}
	}
}

func ServeDashboard(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "index.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Admin dashboard template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServePostsPage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "posts.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Posts page template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeNewPostPage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "post_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "New post form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeEditPostPage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "post_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Edit post form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServePortfolioPage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "portfolio_list.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Portfolio page template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeNewPortfolioPage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "portfolio_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "New portfolio form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeEditPortfolioPage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "portfolio_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Edit portfolio form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServePagesPage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "page_list.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Pages page template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeNewPagePage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "page_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "New page form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeEditPagePage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "page_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Page edit template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeSettingsPage(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(AdminFilesPath, "settings.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Settings template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}
