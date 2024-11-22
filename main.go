package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/yuin/goldmark"
	// highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func templateFromBase(names ...string) (*template.Template, error) {
	templates := append(names, "./templates/base.html", "./templates/footer.html", "./templates/navbar.html")

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return tmpl, nil
}

func serveAboutMe(w http.ResponseWriter, r *http.Request) {
	tmpl, err := templateFromBase("./pages/home/home.html", "./pages/home/projects.html",
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
	files, err := os.ReadDir("./posts")
	if err != nil {
		log.Fatal("Cannot open posts folder")
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	// Extract files metadata
	var posts []Metadata
	for _, f := range files {
		content, err := os.ReadFile(path.Join("./posts", f.Name()))
		if err != nil {
			continue
		}

		// extract file metadata
		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
			panic(err)
		}
		metaData := meta.Get(context)

		posts = append(posts, Metadata{
			Title:       metaData["Title"],
			Date:        metaData["Date"],
			Description: metaData["Description"],
			Image:       metaData["Image"],
		})
	}

	tmp, err := templateFromBase("./pages/posts/list.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	err = tmp.ExecuteTemplate(w, "base.html", &TemplateData{Posts: posts})
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
	}
}

func main() {
	// route for serving static files
	assets_path := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets_path))

	// pages
	http.HandleFunc("/", serveAboutMe)
	http.HandleFunc("/posts/", servePostsList)
	// http.HandleFunc("/posts/{id}/", servePostDetail)

	log.Print("Starting server...\n")
	log.Print("Serving in port 8000\n")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Cannot initialize server on port 8000: %s", err.Error())
	}
}

type TemplateData struct {
	Posts []Metadata
}

// NOTE: for now is not needed to cast it to some specific value and curate the data
type Metadata struct {
	Title       interface{}
	Date        interface{}
	Description interface{}
	Image       interface{}
}

// markdown2 := goldmark.New(
// 	goldmark.WithExtensions(
// 		meta.Meta,
// 		highlighting.NewHighlighting(
// 			highlighting.WithStyle("monokai"),
// 			highlighting.WithFormatOptions(
// 				chromahtml.WithLineNumbers(true),
// 			),
// 		),
// 	),
// )
