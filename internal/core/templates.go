package core

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/main.tex.tmpl
var mainTexTemplate string

//go:embed templates/ieee.tex.tmpl
var ieeeTemplate string

//go:embed templates/acm.tex.tmpl
var acmTemplate string

//go:embed templates/nature.tex.tmpl
var natureTemplate string

//go:embed templates/springer.tex.tmpl
var springerTemplate string

//go:embed templates/cvpr.tex.tmpl
var cvprTemplate string

//go:embed templates/refs.bib.tmpl
var refsBibTemplate string

//go:embed templates/gitignore.tmpl
var gitignoreTemplate string

// GetAvailableTemplates returns a list of available template names.
func GetAvailableTemplates() []string {
	templates := []string{"article", "ieee", "acm", "nature", "springer", "cvpr"}

	// Add custom templates from ~/.config/navtex/templates/
	home, err := os.UserHomeDir()
	if err == nil {
		customDir := filepath.Join(home, ".config", "navtex", "templates")
		files, err := os.ReadDir(customDir)
		if err == nil {
			for _, f := range files {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".tex.tmpl") {
					name := strings.TrimSuffix(f.Name(), ".tex.tmpl")
					templates = append(templates, name)
				}
			}
		}
	}

	return templates
}

// loadTemplate loads a template by name, checking custom templates first.
func loadTemplate(name string) (string, error) {
	// Check custom templates first
	home, err := os.UserHomeDir()
	if err == nil {
		customPath := filepath.Join(home, ".config", "navtex", "templates", name+".tex.tmpl")
		content, err := os.ReadFile(customPath)
		if err == nil {
			return string(content), nil
		}
	}

	// Fallback to embedded templates
	switch strings.ToLower(name) {
	case "article", "main":
		return mainTexTemplate, nil
	case "ieee":
		return ieeeTemplate, nil
	case "acm":
		return acmTemplate, nil
	case "nature":
		return natureTemplate, nil
	case "springer":
		return springerTemplate, nil
	case "cvpr":
		return cvprTemplate, nil
	default:
		// If not found as a template, treat it as a document class using the default article template
		return strings.ReplaceAll(mainTexTemplate, "{{DOCCLASS}}", name), nil
	}
}

// CreateProject scaffolds a new LaTeX project in the given directory.
func CreateProject(root, title, author, templateName string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	if templateName == "" {
		templateName = "article"
	}

	// Create images directory
	imagesDir := filepath.Join(absRoot, "images")
	if err := os.MkdirAll(imagesDir, 0o755); err != nil {
		return fmt.Errorf("creating images directory: %w", err)
	}

	// Load template
	templateContent, err := loadTemplate(templateName)
	if err != nil {
		return fmt.Errorf("loading template %s: %w", templateName, err)
	}

	// Write main.tex
	mainContent := strings.ReplaceAll(templateContent, "{{TITLE}}", title)
	mainContent = strings.ReplaceAll(mainContent, "{{AUTHOR}}", author)
	// If the template has {{DOCCLASS}}, replace it with article if templateName is not "article"
	// and matches an embedded template (some templates have hardcoded classes).
	// For the generic "article" template, we use the provided name as the class.
	if strings.Contains(templateContent, "{{DOCCLASS}}") {
		class := templateName
		if templateName == "article" {
			class = "article"
		}
		mainContent = strings.ReplaceAll(mainContent, "{{DOCCLASS}}", class)
	}
	mainContent = strings.ReplaceAll(mainContent, "{{DATE}}", "\\today")

	if err := writeIfNotExists(filepath.Join(absRoot, "main.tex"), mainContent); err != nil {
		return err
	}

	// Write refs.bib
	if err := writeIfNotExists(filepath.Join(absRoot, "refs.bib"), refsBibTemplate); err != nil {
		return err
	}

	// Write .gitignore
	if err := writeIfNotExists(filepath.Join(absRoot, ".gitignore"), gitignoreTemplate); err != nil {
		return err
	}

	// Create a .gitkeep in images/ so git tracks the empty directory
	gitkeep := filepath.Join(imagesDir, ".gitkeep")
	if err := writeIfNotExists(gitkeep, ""); err != nil {
		return err
	}

	return nil
}

// writeIfNotExists writes content to a file only if the file doesn't already exist.
func writeIfNotExists(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil // File already exists, skip
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", filepath.Base(path), err)
	}
	return nil
}
