package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/elias-gill/portfolio/logger"
)

const sitemapURL = "https://portfolio-elias-gill.fly.dev/posts/sitemap.xml/"

var (
	pingURL = "https://www.google.com/ping?sitemap=" + url.QueryEscape(sitemapURL)

	// Cache solo de base + parciales
	baseTemplates     *template.Template
	baseTemplatesOnce sync.Once

	// I dont know a better name, but here is where the source code is stored, so we can load
	// templates and assets
	resourcesPath        = "."
	assets_path          string
	posts_media_path     string
	templates_path       string
	templates_pages_path string
)

func RegisterRoutes() {
	secret = logger.GetEnvVarAndLog("WEBHOOK_SECRET")
	blogPath = logger.GetEnvVarAndLog("BLOG_PATH")

	envResourcesPath := logger.GetEnvVarAndLog("RESOURCES_PATH")
	if envResourcesPath != "" {
		resourcesPath = envResourcesPath
	}

	posts_media_path = path.Join(blogPath, "media")
	assets_path = path.Join(resourcesPath, "assets")
	templates_path = path.Join(resourcesPath, "templates")
	templates_pages_path = path.Join(templates_path, "pages")

	// route for serving static files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assets_path))))

	// route for serving posts media and attachments
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(posts_media_path))))

	// pages
	http.HandleFunc("/webhook", handleWebhook)
	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/posts/", serveBlogIndex)
	http.HandleFunc("/posts/{post}/", serveBlogpostDetail)

	// search engines indexing
	http.HandleFunc("/robots.txt", handleRobots)
}

// ============================================
// 			Route handlers
// ============================================

func handleRobots(w http.ResponseWriter, r *http.Request) {
	robotsTxt := `User-agent: *
	Allow: /
	`
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(robotsTxt))
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplates(
		[]string{
			path.Join(templates_pages_path, "home", "home.html"),
			path.Join(templates_pages_path, "home", "projects.html"),
			path.Join(templates_pages_path, "home", "tecnologias.html"),
		},
		nil,
		w,
	)
}

func serveBlogIndex(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(blogPath)
	if err != nil {
		respondError(w, "Cannot open posts folder", err, http.StatusInternalServerError)
		return
	}

	var posts []Metadata
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".md") {
			continue
		}

		data, err := extractMetaFromDirEntry(f)
		if err != nil {
			logger.Warn("Failed to extract metadata", "file", f.Name(), "error", err)
			continue
		}
		posts = append(posts, *data)
	}

	renderTemplates(
		[]string{path.Join(templates_pages_path, "posts", "blogIndex.html")},
		map[string]any{"posts": posts},
		w,
	)
}

func serveBlogpostDetail(w http.ResponseWriter, r *http.Request) {
	postFileName := r.PathValue("post")
	if postFileName == "" || strings.Contains(postFileName, "..") {
		respondError(w, "Invalid post path", nil, http.StatusBadRequest)
		return
	}

	fileContent, err := os.ReadFile(path.Join(blogPath, postFileName))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// NOTE: para poder generar sitemaps "dinamicos". Listamos el sitemap que ya esta en la
	// blog path. No me cuesta nada agregar al sitemap una entrada nueva de forma manual.
	if postFileName == "sitemap.xml" {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fileContent))
		return
	}

	var parsedMarkdown bytes.Buffer
	if err := parseFile(fileContent, &parsedMarkdown); err != nil {
		respondError(w, "Error loading markdown file", err, http.StatusInternalServerError)
		return
	}

	fileMetadata, err := extractMetadataFromFilePath(postFileName)
	if err != nil {
		respondError(w, "Error loading file metadata", err, http.StatusInternalServerError)
		return
	}

	post := Post{
		Content:  template.HTML(parsedMarkdown.String()),
		MetaData: fileMetadata,
	}

	renderTemplates(
		[]string{path.Join(templates_pages_path, "posts", "blogpostDetail.html")},
		&post,
		w,
	)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Verify the secret
	if !verifySecret(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Pull the latest changes from the remote repository
	err := gitPull()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to pull: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Repository updated successfully"))
}

// ============================================
// 			Template rendering
// ============================================

func initBaseTemplates() {
	baseTemplatesOnce.Do(func() {
		// Registrar las helpers-funcs para las plantillas
		funcMap := template.FuncMap{
			"fileExists": fileExists,
		}

		baseTemplates = template.New("").Funcs(funcMap)

		files := []string{
			path.Join(templates_path, "base.html"),
			path.Join(templates_path, "navbar.html"),
			path.Join(templates_path, "footer.html"),
		}

		var err error
		baseTemplates, err = baseTemplates.ParseFiles(files...)
		if err != nil {
			logger.Error("Error parsing base templates", "error", err)
		}
	})
}

func renderTemplates(pageFiles []string, data any, w http.ResponseWriter) {
	initBaseTemplates()

	// Clone the base template
	tmpl, err := baseTemplates.Clone()
	if err != nil {
		logger.Error("Error cloning base template", "error", err)
		respondError(w, "Error cloning base template", err, http.StatusInternalServerError)
		return
	}

	// Parse the required template pages
	tmpl, err = tmpl.ParseFiles(pageFiles...)
	if err != nil {
		logger.Error("Error parsing page templates",
			"files", strings.Join(pageFiles, ","),
			"error", err)
		respondError(w, fmt.Sprintf("Error parsing page templates [%s]", strings.Join(pageFiles, ",")), err, http.StatusInternalServerError)
		return
	}

	// Render the required page, using the base template and with the given data
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		logger.Error("Template execution error",
			"files", strings.Join(pageFiles, ","),
			"error", err)
		respondError(w, fmt.Sprintf("Template execution error for [%s]", strings.Join(pageFiles, ",")), err, http.StatusInternalServerError)
	}
}

// ============================================
// 		Utilities and helper functions
// ============================================

// Helper para centralizar errores HTTP
func respondError(w http.ResponseWriter, msg string, err error, code int) {
	if err != nil {
		logger.Error(msg, "error", err)
	} else {
		logger.Warn(msg)
	}
	http.Error(w, msg, code)
}

// FileExists helper
func fileExists(dir string) bool {
	if dir == "" {
		return false
	}
	_, err := os.Stat(path.Join(posts_media_path, dir))
	return err == nil
}

func verifySecret(r *http.Request) bool {
	// Get the signature from the request header
	signature := r.Header.Get("X-Hub-Signature")
	if signature == "" {
		return false
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return false
	}
	defer r.Body.Close()

	// Compute the HMAC signature
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(body)
	expectedSignature := "sha1=" + hex.EncodeToString(mac.Sum(nil))

	// Compare the signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// If the repository does not exist at the specified path, it clones the repository.
func gitPull() error {
	// Pull the latest changes if the repository exists
	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = blogPath

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, out.String())
	}

	logger.Info("Updated repo succesfully")

	if err := pingGoogle(); err != nil {
		logger.Error("Error updating blog sitemap", "error", err)
	}

	return nil
}

func pingGoogle() error {
	resp, err := http.Get(pingURL)
	if err != nil {
		return fmt.Errorf("Error haciendo ping a Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Google devolvió status %d", resp.StatusCode)
	}

	return nil
}
