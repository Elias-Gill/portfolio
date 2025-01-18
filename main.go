package main

import (
	"log"
	"net/http"
)

const (
	secret        = "your-secret-token"        // TODO: Replace with your webhook secret
	repoPath      = "/path/to/your/local/repo" // TODO: Replace with your local repo path
	remoteRepoURL = "https://github.com/elias-gill/blog"
	port          = ":8000"
)

func main() {
	// route for serving static files
	assets_path := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets_path))

	// route for serving posts media and attachments
	posts_media_path := http.FileServer(http.Dir("./posts/media/"))
	http.Handle("/media/", http.StripPrefix("/media/", posts_media_path))

	// pages
	http.HandleFunc("/webhook", handleWebhook)
	http.HandleFunc("/", serveAboutMe)
	http.HandleFunc("/posts/", servePostsList)
	http.HandleFunc("/posts/{id}/", servePostDetail)

	log.Print("Starting server...\n")
	log.Printf("Serving in port %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Cannot initialize server on port %s %s", port, err.Error())
	}
}
