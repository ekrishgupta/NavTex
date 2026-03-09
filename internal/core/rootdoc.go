package core

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// magicRootRe matches lines like: % !TEX root = ../main.tex
var magicRootRe = regexp.MustCompile(`(?i)^\s*%\s*!TEX\s+root\s*=\s*(.+)`)

// ResolveRootDocument attempts to find the main .tex file for a given file.
// 1. It checks the first 50 lines for a magic comment (% !TEX root = ...).
// 2. It checks if the current file has a \documentclass.
// 3. Defaults to scanning the project root directory for a .tex file with \documentclass.
func ResolveRootDocument(texPath, projectRootDir string) (string, error) {
	file, err := os.Open(texPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	hasDocumentClass := false
	lineCount := 0

	// Scan the first 50 lines to detect magic comments or \documentclass
	for scanner.Scan() {
		line := scanner.Text()

		// 1. Check for magic comment
		matches := magicRootRe.FindStringSubmatch(line)
		if len(matches) > 1 {
			magicPath := strings.TrimSpace(matches[1])
			// Resolve relative to the current file's directory
			resolvedPath, err := filepath.Abs(filepath.Join(filepath.Dir(texPath), magicPath))
			if err == nil {
				return resolvedPath, nil
			}
		}

		// 2. Check for preamble
		if strings.Contains(line, `\documentclass`) {
			hasDocumentClass = true
			break
		}

		lineCount++
		if lineCount > 100 {
			break
		}
	}

	if !hasDocumentClass {
		data, readErr := os.ReadFile(texPath)
		if readErr == nil {
			if strings.Contains(string(data), `\documentclass`) {
				hasDocumentClass = true
			}
		}
	}

	if hasDocumentClass {
		return texPath, nil
	}

	// 3. Fallback: Scan project root directory for a master file
	masterFile := ""
	filepath.WalkDir(projectRootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && path != projectRootDir {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.ToLower(filepath.Ext(path)) == ".tex" {
			data, err := os.ReadFile(path)
			if err == nil && strings.Contains(string(data), `\documentclass`) {
				masterFile = path
				return filepath.SkipAll
			}
		}
		return nil
	})

	if masterFile != "" {
		return masterFile, nil
	}

	// If everything fails, just return the original path.
	return texPath, nil
}
