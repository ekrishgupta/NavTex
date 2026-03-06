package core

import (
	"os"
	"testing"
)

func TestPreviewPurge(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some files
	os.WriteFile(tmpDir+"/main.tex", []byte(""), 0o644)
	os.WriteFile(tmpDir+"/main.aux", []byte(""), 0o644)
	os.WriteFile(tmpDir+"/main.log", []byte(""), 0o644)
	os.WriteFile(tmpDir+"/image.png", []byte(""), 0o644)

	entries, err := PreviewPurge(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Should identify .aux and .log but not .tex or .png
	if len(entries) != 2 {
		t.Errorf("Expected 2 preview entries, got %d: %v", len(entries), entries)
	}

	// Verify files still exist
	if _, err := os.Stat(tmpDir + "/main.aux"); os.IsNotExist(err) {
		t.Error("Preview should not delete files")
	}
}

func TestPurge(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(tmpDir+"/main.tex", []byte(""), 0o644)
	os.WriteFile(tmpDir+"/main.aux", []byte(""), 0o644)

	removed, err := Purge(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(removed) != 1 || removed[0] != "main.aux" {
		t.Errorf("Expected [main.aux] to be removed, got %v", removed)
	}

	if _, err := os.Stat(tmpDir + "/main.aux"); !os.IsNotExist(err) {
		t.Error("File main.aux should have been deleted")
	}
	if _, err := os.Stat(tmpDir + "/main.tex"); os.IsNotExist(err) {
		t.Error("File main.tex should NOT have been deleted")
	}
}
