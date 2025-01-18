package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
)

// Generates an extended template using the base.html template
// to ensure a consistent layout and structure across all pages.
func renderTemplates(templates ...string) (*template.Template, error) {
	// Add helper functions
	funcMap := template.FuncMap{
		"fileExists": func(path string) bool {
			if path == "" {
				return false
			}
			_, err := os.Stat(path)
			return !errors.Is(err, fs.ErrExist)
		},
	}

	// Create a new template with the FuncMap
	tmpl := template.New("").Funcs(funcMap)

	// Append the base templates to the list of templates
	tmpls := append(templates, "./templates/base.html", "./templates/footer.html", "./templates/navbar.html")

	// Parse all templates
	tmpl, err := tmpl.ParseFiles(tmpls...)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return tmpl, nil
}

func serveAboutMe(w http.ResponseWriter, r *http.Request) {
	tmpl, err := renderTemplates("./pages/home/home.html", "./pages/home/projects.html",
		"./pages/home/tecnologias.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err = tmpl.ExecuteTemplate(w, "base.html", nil); err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
	}
}

func servePostsList(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(repoPath)
	if err != nil {
		log.Fatal("Cannot open posts folder")
	}

	var posts []Metadata
	for _, f := range files {
		data, err := extractMetaFromDirEntry(f)
		if err != nil {
			continue
		}
		posts = append(posts, *data)
	}

	tmp, err := renderTemplates("./pages/posts/list.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	err = tmp.ExecuteTemplate(w, "base.html", map[string]interface{}{"posts": posts})
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
}

func servePostDetail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("post")
	file, err := os.ReadFile(path.Join(repoPath, id))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("vim"),
			),
		),
	)

	// markdown parsing
	var content bytes.Buffer
	if err = markdown.Convert(file, &content); err != nil {
		http.Error(w, "Error laoding file", http.StatusInternalServerError)
		log.Println(err.Error())
	}

	data, err := extractPostMetadata(id)
	if err = markdown.Convert(file, &content); err != nil {
		http.Error(w, "Error laoding file", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	// template generation
	tmpl, err := renderTemplates("./pages/posts/detail.html")
	if err != nil {
		http.Error(w, "Error laoding template", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	post := Post{
		Content: template.HTML(content.String()),
		Meta:    data,
	}
	tmpl.ExecuteTemplate(w, "base.html", &post)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Verify the secret (optional but recommended for security)
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
