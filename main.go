package main

import (
	"net/http"

	"github.com/elias-gill/portfolio/logger"
)

var (
	secret   string
	blogPath string
	port     = "8000"
)

func main() {
	envPort := logger.GetEnvVarAndLog("PORT")
	if envPort != "" {
		port = envPort
	}

	// Load environment variables
	secret = logger.GetEnvVarAndLog("WEBHOOK_SECRET")
	blogPath = logger.GetEnvVarAndLog("BLOG_PATH")

	RegisterRoutes()

	logger.Info("Starting server...", "port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("Cannot initialize server on port", "port", port, "stack_error", err.Error())
	}
}
