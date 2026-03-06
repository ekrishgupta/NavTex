package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// auxExtensions lists all auxiliary file extensions to purge.
var auxExtensions = []string{
	".aux", ".log", ".nav", ".out", ".snm", ".toc",
	".fls", ".fdb_latexmk", ".synctex.gz",
	".bbl", ".blg", ".lof", ".lot",
	".idx", ".ind", ".ilg",
	".vrb", ".xdv", ".bcf", ".run.xml",
	".dvi", ".ps",
}

// PreviewPurge returns the list of auxiliary files that would be deleted.
func PreviewPurge(root string) ([]string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	var targets []string

	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && path != absRoot {
				return filepath.SkipDir
			}
			return nil
		}

		if isAuxFile(info.Name()) {
			rel, _ := filepath.Rel(absRoot, path)
			targets = append(targets, rel)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return targets, nil
}

// Purge deletes all auxiliary files from the given directory and returns
// the list of files that were removed.
func Purge(root string) ([]string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	var removed []string

	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && path != absRoot {
				return filepath.SkipDir
			}
			return nil
		}

		if isAuxFile(info.Name()) {
			if err := os.Remove(path); err != nil {
				return nil // Best-effort: skip files we can't delete
			}
			rel, _ := filepath.Rel(absRoot, path)
			removed = append(removed, rel)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return removed, nil
}

// isAuxFile checks if a filename has an auxiliary extension.
func isAuxFile(name string) bool {
	lower := strings.ToLower(name)
	for _, ext := range auxExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}
