package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateProject_Basic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = CreateProject(tmpDir, "My Paper", "John Doe", "report")
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	// Check if files exist
	expectedFiles := []string{
		"main.tex",
		"refs.bib",
		".gitignore",
		"images",
		"images/.gitkeep",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file/dir %s does not exist", f)
		}
	}
}

func TestCreateProject_SkipExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	texPath := filepath.Join(tmpDir, "main.tex")
	originalContent := []byte("original")
	os.WriteFile(texPath, originalContent, 0o644)

	err = CreateProject(tmpDir, "New Paper", "Author", "article")
	if err != nil {
		t.Fatal(err)
	}

	content, _ := os.ReadFile(texPath)
	if string(content) != string(originalContent) {
		t.Error("CreateProject should not overwrite existing files")
	}
}
