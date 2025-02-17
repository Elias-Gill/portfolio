package main

import (
	"log"
	"net/http"
	"path"
)

const (
	remoteRepoURL = "https://github.com/elias-gill/blog"
)

var (
	secret   string
	repoPath string
	mainPath = "."
	port     = "8000"
)

func main() {
	// Load environment variables
	secret = getEnvAndLog("WEBHOOK_SECRET")
	repoPath = getEnvAndLog("BLOG_PATH")

	envPort := getEnvAndLog("PORT")
	if envPort != "" {
		port = envPort
	}

	envMPath := getEnvAndLog("MAIN_PATH")
	if envPort != "" {
		mainPath = envMPath
	}

	// route for serving static files
	assets_path := http.FileServer(http.Dir(path.Join(mainPath, "/assets")))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets_path))

	// route for serving posts media and attachments
	posts_media_path := http.FileServer(http.Dir(path.Join(repoPath, "/media/")))
	http.Handle("/media/", http.StripPrefix("/media/", posts_media_path))

	// pages
	http.HandleFunc("/webhook", handleWebhook)
	http.HandleFunc("/", serveAboutMe)
	http.HandleFunc("/posts/", serveBlog)
	http.HandleFunc("/posts/{post}/", servePostDetail)

	// search engines indexing
	http.HandleFunc("/robots.txt", robotsHandler)

	log.Print("Starting server...\n")
	log.Printf("Serving in port %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Cannot initialize server on port %s %s", port, err.Error())
	}
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	robotsTxt := `User-agent: *
	Allow: /
	`
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(robotsTxt))
}
