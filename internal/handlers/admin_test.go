package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeDashboard(t *testing.T) {
	// Set ADMIN_FILES_PATH to project directory for testing
	AdminFilesPath = "../../../personal-blog-generator/admin-files"

	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	w := httptest.NewRecorder()

	ServeDashboard(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Admin Dashboard") {
		t.Errorf("Expected body to contain 'Admin Dashboard', got: %s", body)
	}
}
