package handlers

import (
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

var TemplatePath string
var OutputPath string
var AdminFS fs.FS

// AdminPageData holds data for admin page templates
type AdminPageData struct {
	Title        string
	ActiveNav    string
	ExtraHead    template.HTML
	Scripts      template.HTML
	Content      template.HTML
	TemplatePath string
	OutputPath   string
}

func init() {
	// Get TEMPLATE_PATH and OUTPUT_PATH for dashboard display
	TemplatePath = os.Getenv("TEMPLATE_PATH")
	if TemplatePath == "" {
		homeDir, _ := os.UserHomeDir()
		TemplatePath = filepath.Join(homeDir, ".personal-blog-generator", "templates")
	}

	OutputPath = os.Getenv("OUTPUT_PATH")
	if OutputPath == "" {
		homeDir, _ := os.UserHomeDir()
		OutputPath = filepath.Join(homeDir, ".personal-blog-generator", "html-outputs")
	}
}

// renderAdminPage renders an admin page using header.html, content, and footer.html
func renderAdminPage(w http.ResponseWriter, data AdminPageData) error {
	// Read header and footer templates
	headerBytes, err := fs.ReadFile(AdminFS, "header.html")
	if err != nil {
		return err
	}
	footerBytes, err := fs.ReadFile(AdminFS, "footer.html")
	if err != nil {
		return err
	}

	// Combine templates
	fullTemplate := string(headerBytes) + string(data.Content) + string(footerBytes)

	tmpl, err := template.New("admin").Parse(fullTemplate)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.Execute(w, data)
}

// readContentFile reads the content portion of an admin HTML file (stripping header/footer)
func readContentFile(filename string) (template.HTML, error) {
	content, err := fs.ReadFile(AdminFS, "content/"+filename)
	if err != nil {
		return "", err
	}
	return template.HTML(content), nil
}

func ServeDashboard(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("dashboard.html")
	if err != nil {
		http.Error(w, "Admin dashboard template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:        "Admin Dashboard",
		ActiveNav:    "dashboard",
		Content:      content,
		TemplatePath: TemplatePath,
		OutputPath:   OutputPath,
		Scripts: template.HTML(`<script>
        document.getElementById('publish-site-btn').addEventListener('click', async function() {
            const btn = this;
            const originalText = btn.innerHTML;
            btn.innerHTML = '<span class="material-symbols-outlined">refresh</span><span>Publishing...</span>';
            btn.disabled = true;

            try {
                const response = await fetch('/api/publish', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                const result = await response.json();

                if (response.ok) {
                    alert('Site published successfully! Generated ' + result.count + ' post pages.');
                } else {
                    alert('Publish failed: ' + result.error);
                }
            } catch (error) {
                console.error('Publish error:', error);
                alert('Network error. Please try again.');
            } finally {
                btn.innerHTML = originalText;
                btn.disabled = false;
            }
        });
    </script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
	}
}

func ServePostsPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("posts.html")
	if err != nil {
		http.Error(w, "Posts page template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "Blog Posts",
		ActiveNav: "posts",
		Content:   content,
		Scripts:   template.HTML(`<script src="/admin/js/posts.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render posts page", http.StatusInternalServerError)
	}
}

func ServeNewPostPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("post_form.html")
	if err != nil {
		http.Error(w, "New post form template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "New Post",
		ActiveNav: "posts",
		Content:   content,
		ExtraHead: template.HTML(`<link rel="stylesheet" href="/admin/vendor/easymde.min.css">`),
		Scripts:   template.HTML(`<script src="/admin/vendor/easymde.min.js"></script><script src="/admin/js/post_form.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render new post page", http.StatusInternalServerError)
	}
}

func ServeEditPostPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("post_form.html")
	if err != nil {
		http.Error(w, "Edit post form template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "Edit Post",
		ActiveNav: "posts",
		Content:   content,
		ExtraHead: template.HTML(`<link rel="stylesheet" href="/admin/vendor/easymde.min.css">`),
		Scripts:   template.HTML(`<script src="/admin/vendor/easymde.min.js"></script><script src="/admin/js/post_form.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render edit post page", http.StatusInternalServerError)
	}
}

func ServePortfolioPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("portfolio_list.html")
	if err != nil {
		http.Error(w, "Portfolio page template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "Portfolio",
		ActiveNav: "portfolio",
		Content:   content,
		Scripts:   template.HTML(`<script src="/admin/js/portfolio_list.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render portfolio page", http.StatusInternalServerError)
	}
}

func ServeNewPortfolioPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("portfolio_form.html")
	if err != nil {
		http.Error(w, "New portfolio form template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "New Portfolio Item",
		ActiveNav: "portfolio",
		Content:   content,
		ExtraHead: template.HTML(`<link rel="stylesheet" href="/admin/vendor/easymde.min.css">`),
		Scripts:   template.HTML(`<script src="/admin/vendor/easymde.min.js"></script><script src="/admin/js/portfolio_form.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render new portfolio page", http.StatusInternalServerError)
	}
}

func ServeEditPortfolioPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("portfolio_form.html")
	if err != nil {
		http.Error(w, "Edit portfolio form template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "Edit Portfolio Item",
		ActiveNav: "portfolio",
		Content:   content,
		ExtraHead: template.HTML(`<link rel="stylesheet" href="/admin/vendor/easymde.min.css">`),
		Scripts:   template.HTML(`<script src="/admin/vendor/easymde.min.js"></script><script src="/admin/js/portfolio_form.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render edit portfolio page", http.StatusInternalServerError)
	}
}

func ServePagesPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("page_list.html")
	if err != nil {
		http.Error(w, "Pages page template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "Pages",
		ActiveNav: "pages",
		Content:   content,
		Scripts:   template.HTML(`<script src="/admin/js/page_list.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render pages page", http.StatusInternalServerError)
	}
}

func ServeNewPagePage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("page_form.html")
	if err != nil {
		http.Error(w, "New page form template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "New Page",
		ActiveNav: "pages",
		Content:   content,
		ExtraHead: template.HTML(`<link rel="stylesheet" href="/admin/vendor/easymde.min.css">`),
		Scripts:   template.HTML(`<script src="/admin/vendor/easymde.min.js"></script><script src="/admin/js/page_form.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render new page page", http.StatusInternalServerError)
	}
}

func ServeEditPagePage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("page_form.html")
	if err != nil {
		http.Error(w, "Page edit template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "Edit Page",
		ActiveNav: "pages",
		Content:   content,
		ExtraHead: template.HTML(`<link rel="stylesheet" href="/admin/vendor/easymde.min.css">`),
		Scripts:   template.HTML(`<script src="/admin/vendor/easymde.min.js"></script><script src="/admin/js/page_form.js"></script>`),
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render edit page page", http.StatusInternalServerError)
	}
}

func ServeSettingsPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("settings.html")
	if err != nil {
		http.Error(w, "Settings template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "Settings",
		ActiveNav: "settings",
		Content:   content,
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render settings page", http.StatusInternalServerError)
	}
}

func ServeTemplatesPage(w http.ResponseWriter, r *http.Request) {
	content, err := readContentFile("templates.html")
	if err != nil {
		http.Error(w, "Templates template not found", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Title:     "Templates",
		ActiveNav: "templates",
		Content:   content,
	}

	if err := renderAdminPage(w, data); err != nil {
		http.Error(w, "Failed to render templates page", http.StatusInternalServerError)
	}
}
