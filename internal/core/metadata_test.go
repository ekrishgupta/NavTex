package core

import (
	"os"
	"testing"
)

func TestTexMetadata_Basic(t *testing.T) {
	content := `
\documentclass[12pt]{article}
\usepackage{amsmath}
\usepackage{graphicx}
\title{Test Document}
\author{Author Name}
\begin{document}
Hello world.
\end{document}
`
	tmpFile, err := os.CreateTemp("", "test-*.tex")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	os.WriteFile(tmpFile.Name(), []byte(content), 0o644)

	meta, err := TexMetadata(tmpFile.Name())
	if err != nil {
		t.Fatalf("TexMetadata failed: %v", err)
	}

	if meta.Title != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", meta.Title)
	}
	if meta.Author != "Author Name" {
		t.Errorf("Expected author 'Author Name', got '%s'", meta.Author)
	}
	if meta.DocumentClass != "article" {
		t.Errorf("Expected class 'article', got '%s'", meta.DocumentClass)
	}
}

func TestTexMetadata_WordCount(t *testing.T) {
	content := `
\documentclass{article}
\begin{document}
One two three four five.
\end{document}
`
	tmpFile, err := os.CreateTemp("", "test-*.tex")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	os.WriteFile(tmpFile.Name(), []byte(content), 0o644)

	meta, err := TexMetadata(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if meta.WordCount != 5 {
		t.Errorf("Expected 5 words, got %d", meta.WordCount)
	}
}
