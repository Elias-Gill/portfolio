package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"os"
	"path"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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
	Content  template.HTML
	MetaData *Metadata
}

var markdownMetadataExtractor = goldmark.New(
	goldmark.WithExtensions(
		meta.Meta,
	),
)

// Goldmark parser final
var markdownParser = goldmark.New(
	goldmark.WithRendererOptions(html.WithUnsafe()),
	goldmark.WithExtensions(
		meta.Meta,
		extension.Table,
		highlighting.NewHighlighting(highlighting.WithStyle("base16-snazzy")),
	),
)

func extractMetadataFromFilePath(filePath string) (*Metadata, error) {
	content, err := os.ReadFile(path.Join(blogPath, filePath))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	context := parser.NewContext()
	err = markdownMetadataExtractor.Convert(content, &buf, parser.WithContext(context))
	if err != nil {
		return nil, err
	}

	metaData := meta.Get(context)
	var img string
	if metaData["Image"] != nil {
		img = metaData["Image"].(string)
	}

	return &Metadata{
		Id:          filePath,
		Title:       metaData["Title"],
		Date:        metaData["Date"],
		Description: metaData["Description"],
		Image:       img,
	}, nil
}

// version for DirEntry just calls the unified function
func extractMetaFromDirEntry(file fs.DirEntry) (*Metadata, error) {
	return extractMetadataFromFilePath(file.Name())
}

func parseFile(body []byte, ouputBuffer *bytes.Buffer) error {
	return markdownParser.Convert(body, ouputBuffer)
}
