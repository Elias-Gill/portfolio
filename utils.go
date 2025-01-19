package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// NOTE: for now is not needed to cast the metadata to specific data types
type Metadata struct {
	Title       interface{}
	Date        interface{}
	Description interface{}
	Image       string
	Id          interface{}
}

type Post struct {
	Content template.HTML
	Meta    *Metadata
}

// Like extractPostMetadata, but takes in a fs.DirEntry, so the steps to retrieve the metadata
// is a little bit different.
func extractMetaFromDirEntry(file fs.DirEntry) (*Metadata, error) {
	content, err := os.ReadFile(path.Join(repoPath, file.Name()))
	if err != nil {
		return nil, err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	// extract file metadata
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
		log.Printf("Error parsing metadata: %s\n", err.Error())
		return nil, err
	}
	metaData := meta.Get(context)

	var img string
	if metaData["Image"] != nil {
		img = metaData["Image"].(string)
	}

	return &Metadata{
		Id:          file.Name(),
		Title:       metaData["Title"],
		Date:        metaData["Date"],
		Description: metaData["Description"],
		Image:       img,
	}, nil
}

func extractPostMetadata(file string) (*Metadata, error) {
	content, err := os.ReadFile(path.Join(repoPath, file))
	if err != nil {
		return nil, err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	// extract file metadata
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}
	metaData := meta.Get(context)

	var img string
	if metaData["Image"] != nil {
		img = metaData["Image"].(string)
	}

	return &Metadata{
		Image:       img,
		Id:          file,
		Title:       metaData["Title"],
		Date:        metaData["Date"],
		Description: metaData["Description"],
	}, nil
}

// If the repository does not exist at the specified path, it clones the repository.
func gitPull() error {
	// Pull the latest changes if the repository exists
	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = repoPath

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, out.String())
	}

	log.Printf("Updated repo succesfully")
	return nil
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

// getEnvAndLog retrieves the value of an environment variable and logs if it is not set.
func getEnvAndLog(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Environment variable %s is not set", key)
	}
	return value
}
