package handlers

import (
	"fmt"
	"github.com/ariefbayu/personal-blog-generator/internal/templates/admin"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)
var TemplatePath string
var OutputPath string
var DBPath string
var RootPrefix string
var AdminFS fs.FS

// AdminPageData holds data for admin page templates
type AdminPageData struct {
        Title        string
        ActiveNav    string
        ExtraHead    template.HTML
        Scripts      template.HTML
        TemplatePath string
        OutputPath   string
        DBPath       string
        RootPrefix   string
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

        RootPrefix = os.Getenv("ROOT_PREFIX")
        if RootPrefix == "" {
                RootPrefix = "/admin"
        }
        if RootPrefix == "/" {
                RootPrefix = ""
        } else {
                if !strings.HasPrefix(RootPrefix, "/") {
                        RootPrefix = "/" + RootPrefix
                }
                RootPrefix = strings.TrimSuffix(RootPrefix, "/")
        }
}

// mapToTemplData converts local AdminPageData to the one expected by templ templates
func mapToTemplData(data AdminPageData) admin.AdminPageData {
        return admin.AdminPageData{
                Title:        data.Title,
                ActiveNav:    data.ActiveNav,
                ExtraHead:    data.ExtraHead,
                Scripts:      data.Scripts,
                TemplatePath: data.TemplatePath,
                OutputPath:   data.OutputPath,
                DBPath:       data.DBPath,
                RootPrefix:   data.RootPrefix,
        }
}

func ServeDashboard(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:        "Admin Dashboard",
                ActiveNav:    "dashboard",
                TemplatePath: TemplatePath,
                OutputPath:   OutputPath,
                DBPath:       DBPath,
                RootPrefix:   RootPrefix,
                Scripts: template.HTML(fmt.Sprintf(`<script>
        document.getElementById('publish-site-btn').addEventListener('click', async function() {
            const btn = this;
            const originalText = btn.innerHTML;
            btn.innerHTML = '<span class="material-symbols-outlined">refresh</span><span>Publishing...</span>';
            btn.disabled = true;

            try {
                const response = await fetch('%s/api/publish', {
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
    </script>`, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.Dashboard(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
        }
}
func ServePostsPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "Blog Posts",
                ActiveNav:  "posts",
                RootPrefix: RootPrefix,
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/js/posts.js"></script>`, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.Posts(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render posts page", http.StatusInternalServerError)
        }
}

func ServeNewPostPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "New Post",
                ActiveNav:  "posts",
                RootPrefix: RootPrefix,
                ExtraHead:  template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s/vendor/easymde.min.css">`, RootPrefix)),
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/vendor/easymde.min.js"></script><script src="%s/js/post_form.js"></script>`, RootPrefix, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.PostForm(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render new post page", http.StatusInternalServerError)
        }
}

func ServeEditPostPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "Edit Post",
                ActiveNav:  "posts",
                RootPrefix: RootPrefix,
                ExtraHead:  template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s/vendor/easymde.min.css">`, RootPrefix)),
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/vendor/easymde.min.js"></script><script src="%s/js/post_form.js"></script>`, RootPrefix, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.PostForm(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render edit post page", http.StatusInternalServerError)
        }
}

func ServePortfolioPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "Portfolio",
                ActiveNav:  "portfolio",
                RootPrefix: RootPrefix,
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/js/portfolio_list.js"></script>`, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.PortfolioList(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render portfolio page", http.StatusInternalServerError)
        }
}

func ServeNewPortfolioPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "New Portfolio Item",
                ActiveNav:  "portfolio",
                RootPrefix: RootPrefix,
                ExtraHead:  template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s/vendor/easymde.min.css">`, RootPrefix)),
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/vendor/easymde.min.js"></script><script src="%s/js/portfolio_form.js"></script>`, RootPrefix, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.PortfolioForm(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render new portfolio page", http.StatusInternalServerError)
        }
}

func ServeEditPortfolioPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "Edit Portfolio Item",
                ActiveNav:  "portfolio",
                RootPrefix: RootPrefix,
                ExtraHead:  template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s/vendor/easymde.min.css">`, RootPrefix)),
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/vendor/easymde.min.js"></script><script src="%s/js/portfolio_form.js"></script>`, RootPrefix, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.PortfolioForm(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render edit portfolio page", http.StatusInternalServerError)
        }
}

func ServePagesPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "Pages",
                ActiveNav:  "pages",
                RootPrefix: RootPrefix,
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/js/page_list.js"></script>`, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.PageList(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render pages page", http.StatusInternalServerError)
        }
}

func ServeNewPagePage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "New Page",
                ActiveNav:  "pages",
                RootPrefix: RootPrefix,
                ExtraHead:  template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s/vendor/easymde.min.css">`, RootPrefix)),
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/vendor/easymde.min.js"></script><script src="%s/js/page_form.js"></script>`, RootPrefix, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.PageForm(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render new page page", http.StatusInternalServerError)
        }
}

func ServeEditPagePage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "Edit Page",
                ActiveNav:  "pages",
                RootPrefix: RootPrefix,
                ExtraHead:  template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s/vendor/easymde.min.css">`, RootPrefix)),
                Scripts:    template.HTML(fmt.Sprintf(`<script src="%s/vendor/easymde.min.js"></script><script src="%s/js/page_form.js"></script>`, RootPrefix, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.PageForm(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render edit page page", http.StatusInternalServerError)
        }
}

func ServeSettingsPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "Settings",
                ActiveNav:  "settings",
                RootPrefix: RootPrefix,
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.Settings(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render settings page", http.StatusInternalServerError)
        }
}

func ServeTemplatesPage(w http.ResponseWriter, r *http.Request) {
        data := AdminPageData{
                Title:      "Templates",
                ActiveNav:  "templates",
                RootPrefix: RootPrefix,
                ExtraHead: template.HTML(`
                        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/codemirror.min.css">
                        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/theme/material-ocean.min.css">
                `),
                Scripts: template.HTML(fmt.Sprintf(`
                        <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/codemirror.min.js"></script>
                        <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/mode/xml/xml.min.js"></script>
                        <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/mode/javascript/javascript.min.js"></script>
                        <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/mode/css/css.min.js"></script>
                        <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/mode/htmlmixed/htmlmixed.min.js"></script>
                        <script src="https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.65.16/mode/markdown/markdown.min.js"></script>
                        <script src="%s/js/templates.js"></script>
                `, RootPrefix)),
        }

        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := admin.Templates(mapToTemplData(data)).Render(r.Context(), w); err != nil {
                http.Error(w, "Failed to render templates page", http.StatusInternalServerError)
        }
}
