package main

import (
	"html/template"
	"log"
	"net/http"
)

func templateFromBase(w http.ResponseWriter, names ...string) (*template.Template, error) {
	templates := append(names, "./templates/base.html", "./templates/footer.html", "./templates/navbar.html")

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Print(err.Error())
		return nil, err
	}

	return tmpl, nil
}
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := templateFromBase(w, "./templates/home.html", "./templates/projects.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err = tmpl.ExecuteTemplate(w, "base.html", nil); err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
	}
}

func main() {
	// route for serving static files
	assets_path := http.FileServer(http.Dir("./assets"))
	http.HandleFunc("/", serveTemplate)

	http.Handle("/assets/", http.StripPrefix("/assets/", assets_path))

	print("serving in port 8000:\n")
	http.ListenAndServe(":8000", nil)
}
