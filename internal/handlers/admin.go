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
	filePath := filepath.Join(wd, "admin-files", "index.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Admin dashboard template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServePostsPage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "posts.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Posts page template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeNewPostPage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "post_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "New post form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeEditPostPage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "post_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Edit post form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServePortfolioPage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "portfolio_list.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Portfolio page template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeNewPortfolioPage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "portfolio_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "New portfolio form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeEditPortfolioPage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "portfolio_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Edit portfolio form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServePagesPage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "page_list.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Pages page template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeNewPagePage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "page_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "New page form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}

func ServeEditPagePage(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get working directory", http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(wd, "admin-files", "page_form.html")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Edit page form template not found", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filePath)
}