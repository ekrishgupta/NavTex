package core

import (
	"os"
	"testing"
)

func TestScanDirectory_Empty(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	pf, err := ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if pf.Total() != 0 {
		t.Errorf("Expected 0 files, got %d", pf.Total())
	}
}

func TestScanDirectory_Source(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "navtex-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.WriteFile(tmpDir+"/main.tex", []byte(""), 0o644)

	pf, err := ScanDirectory(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(pf.Source) != 1 {
		t.Errorf("Expected 1 source file, got %d", len(pf.Source))
	}
	if pf.Source[0].Name != "main.tex" {
		t.Errorf("Expected main.tex, got %s", pf.Source[0].Name)
	}
}
