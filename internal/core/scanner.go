package core

import (
	"os"
	"path/filepath"
	"strings"
)

// FileCategory represents the role a file plays in a LaTeX project.
type FileCategory int

const (
	CategorySource    FileCategory = iota // .tex files
	CategoryData                          // .bib, .csv, .py
	CategoryAssets                        // .png, .jpg, .jpeg, .pdf (images)
	CategoryAuxiliary                     // .aux, .log, .nav, .out, etc.
	CategoryOutput                        // .pdf matching a .tex basename
)

// FileEntry represents a single classified file.
type FileEntry struct {
	Path     string       // Absolute path
	Name     string       // Basename
	Category FileCategory // Classified role
	Size     int64        // File size in bytes
}

// ProjectFiles holds all classified files in a LaTeX project directory.
type ProjectFiles struct {
	Root      string
	Source    []FileEntry
	Data      []FileEntry
	Assets    []FileEntry
	Auxiliary []FileEntry
	Output    []FileEntry
}

// Total returns the total number of files across all categories.
func (pf *ProjectFiles) Total() int {
	return len(pf.Source) + len(pf.Data) + len(pf.Assets) + len(pf.Auxiliary) + len(pf.Output)
}

// Extension-to-category mapping.
var (
	sourceExts = map[string]bool{
		".tex": true,
	}

	dataExts = map[string]bool{
		".bib": true,
		".csv": true,
		".py":  true,
		".r":   true,
		".lua": true,
	}

	assetExts = map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".svg":  true,
		".eps":  true,
	}

	auxExts = map[string]bool{
		".aux":         true,
		".log":         true,
		".nav":         true,
		".out":         true,
		".snm":         true,
		".toc":         true,
		".fls":         true,
		".fdb_latexmk": true,
		".synctex.gz":  true,
		".bbl":         true,
		".blg":         true,
		".lof":         true,
		".lot":         true,
		".idx":         true,
		".ind":         true,
		".ilg":         true,
		".vrb":         true,
		".xdv":         true,
		".bcf":         true,
		".run.xml":     true,
	}
)

// ScanDirectory walks the given root directory and classifies every file
// into its LaTeX project role.
func ScanDirectory(root string) (*ProjectFiles, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	pf := &ProjectFiles{Root: absRoot}

	// Collect .tex basenames to identify output PDFs.
	texBasenames := make(map[string]bool)

	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't read
		}
		if info.IsDir() {
			// Skip hidden directories
			if strings.HasPrefix(info.Name(), ".") && path != absRoot {
				return filepath.SkipDir
			}
			return nil
		}

		name := info.Name()
		ext := strings.ToLower(filepath.Ext(name))

		// Handle double extensions like .synctex.gz and .run.xml
		doubleExt := ""
		nameNoExt := strings.TrimSuffix(name, ext)
		if secondExt := filepath.Ext(nameNoExt); secondExt != "" {
			doubleExt = strings.ToLower(secondExt + ext)
		}

		entry := FileEntry{
			Path: path,
			Name: name,
			Size: info.Size(),
		}

		switch {
		case sourceExts[ext]:
			entry.Category = CategorySource
			pf.Source = append(pf.Source, entry)
			texBasenames[strings.TrimSuffix(name, ext)] = true
		case dataExts[ext]:
			entry.Category = CategoryData
			pf.Data = append(pf.Data, entry)
		case assetExts[ext]:
			entry.Category = CategoryAssets
			pf.Assets = append(pf.Assets, entry)
		case auxExts[ext] || auxExts[doubleExt]:
			entry.Category = CategoryAuxiliary
			pf.Auxiliary = append(pf.Auxiliary, entry)
		case ext == ".pdf":
			// PDF is "output" if it matches a .tex basename, else an asset
			baseName := strings.TrimSuffix(name, ext)
			if texBasenames[baseName] {
				entry.Category = CategoryOutput
				pf.Output = append(pf.Output, entry)
			} else {
				entry.Category = CategoryAssets
				pf.Assets = append(pf.Assets, entry)
			}
		default:
			// Ignore files that don't fit any category
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Second pass: reclassify PDFs now that we have all tex basenames.
	// (Walk order may have seen the PDF before the .tex file.)
	var newAssets []FileEntry
	for _, a := range pf.Assets {
		ext := strings.ToLower(filepath.Ext(a.Name))
		if ext == ".pdf" {
			baseName := strings.TrimSuffix(a.Name, ext)
			if texBasenames[baseName] {
				a.Category = CategoryOutput
				pf.Output = append(pf.Output, a)
				continue
			}
		}
		newAssets = append(newAssets, a)
	}
	pf.Assets = newAssets

	return pf, nil
}
