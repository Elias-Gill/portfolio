package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"os"
	"path"

	"github.com/elias-gill/portfolio/logger"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// NOTE: for now is not needed to cast the metadata to specific data types
type Metadata struct {
	Title       any
	Date        any
	Description any
	Image       string
	Id          any
}

type Post struct {
	Content template.HTML
	Meta    *Metadata
}

var markdown = goldmark.New(
	goldmark.WithExtensions(
		meta.Meta,
	),
)

// Like extractPostMetadata, but takes in a fs.DirEntry, so the steps to retrieve the metadata
// is a little bit different.
func extractMetaFromDirEntry(file fs.DirEntry) (*Metadata, error) {
	content, err := os.ReadFile(path.Join(blogPath, file.Name()))
	if err != nil {
		return nil, err
	}

	// extract file metadata
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
		logger.Error("Error parsing metadata", "error", err.Error())
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
	content, err := os.ReadFile(path.Join(blogPath, file))
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
