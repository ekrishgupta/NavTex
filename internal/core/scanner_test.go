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
