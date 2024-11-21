package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
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
	// TODO: Parsear los metadatos y sacar los titulos
	posts, err := os.ReadDir("./posts")
	if err != nil {
		log.Fatal("Cannot open posts folder")
	}

	for _, file := range posts {
		w.Write([]byte(file.Name()))
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
