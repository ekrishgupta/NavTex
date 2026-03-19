package ui

import (
	"testing"

	"github.com/ekrishgupta/navtex/internal/latex"
)

func TestFileBrowser_Filtering(t *testing.T) {
	fb := NewFileBrowser(DefaultStyles().Browser.Focused)

	files := &latex.ProjectFiles{
		Source: []latex.FileEntry{
			{Name: "main.tex", Path: "main.tex"},
			{Name: "chapter1.tex", Path: "chapter1.tex"},
		},
	}
	fb.SetFiles(files)

	// Default should have 2 files + 1 header = 3 items
	if len(fb.items) != 3 {
		t.Errorf("Expected 3 items initially, got %d", len(fb.items))
	}

	// Filter for "chap"
	fb.SetFilter("chap")
	if len(fb.items) != 2 { // 1 header + 1 file
		t.Errorf("Expected 2 items after filter, got %d", len(fb.items))
	}

	if fb.items[1].display != "chapter1.tex" {
		t.Errorf("Expected 'chapter1.tex', got %s", fb.items[1].display)
	}
}
