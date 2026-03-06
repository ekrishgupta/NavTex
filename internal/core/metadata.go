package core

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// TexMeta holds metadata extracted from a .tex file.
type TexMeta struct {
	DocumentClass string
	ClassOptions  string
	Packages      []PackageInfo
	Title         string
	Author        string
	WordCount     int
}

// PackageInfo represents a single \usepackage entry.
type PackageInfo struct {
	Name    string
	Options string
}

// BibEntry represents a single bibliography entry.
type BibEntry struct {
	Key     string
	Type    string // article, book, inproceedings, etc.
	Title   string
	Authors string
	Year    string
	Journal string
}

// ImageMeta holds metadata about an image file.
type ImageMeta struct {
	Width  int
	Height int
	Size   int64
	Format string
}

var (
	reDocumentClass = regexp.MustCompile(`\\documentclass\s*(?:\[([^\]]*)\])?\s*\{([^}]+)\}`)
	reUsePackage    = regexp.MustCompile(`\\usepackage\s*(?:\[([^\]]*)\])?\s*\{([^}]+)\}`)
	reTitle         = regexp.MustCompile(`\\title\s*\{([^}]+)\}`)
	reAuthor        = regexp.MustCompile(`\\author\s*\{([^}]+)\}`)
	reComment       = regexp.MustCompile(`(?m)^%.*$|(?:^[^\\]*)(%.*$)`)
	reBibEntry      = regexp.MustCompile(`@(\w+)\s*\{\s*([^,\s]+)\s*,`)
	reBibField      = regexp.MustCompile(`(?i)\s*(title|author|year|journal)\s*=\s*\{([^}]*)\}`)
	reTexCommand    = regexp.MustCompile(`\\[a-zA-Z]+\*?(?:\[[^\]]*\])*(?:\{[^}]*\})*`)
	reBeginEnd      = regexp.MustCompile(`\\(?:begin|end)\{[^}]+\}`)
)

// TexMetadata extracts metadata from a .tex file.
func TexMetadata(path string) (*TexMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading tex file: %w", err)
	}

	content := string(data)
	meta := &TexMeta{}

	// Extract document class
	if m := reDocumentClass.FindStringSubmatch(content); m != nil {
		meta.ClassOptions = m[1]
		meta.DocumentClass = m[2]
	}

	// Extract packages
	for _, m := range reUsePackage.FindAllStringSubmatch(content, -1) {
		// Handle comma-separated package names
		names := strings.Split(m[2], ",")
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name != "" {
				meta.Packages = append(meta.Packages, PackageInfo{
					Name:    name,
					Options: m[1],
				})
			}
		}
	}

	// Extract title and author
	if m := reTitle.FindStringSubmatch(content); m != nil {
		meta.Title = strings.TrimSpace(m[1])
	}
	if m := reAuthor.FindStringSubmatch(content); m != nil {
		meta.Author = strings.TrimSpace(m[1])
	}

	// Word count: strip comments, commands, and environments, then count words
	meta.WordCount = countWords(content)

	return meta, nil
}

// countWords estimates the word count of a LaTeX document's body text.
func countWords(content string) int {
	// Find \begin{document} and \end{document} to focus on body
	beginIdx := strings.Index(content, `\begin{document}`)
	endIdx := strings.Index(content, `\end{document}`)

	if beginIdx >= 0 && endIdx > beginIdx {
		content = content[beginIdx+len(`\begin{document}`) : endIdx]
	}

	// Remove comments (lines starting with %)
	lines := strings.Split(content, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "%") {
			continue
		}
		// Remove inline comments (not escaped)
		if idx := strings.Index(line, "%"); idx > 0 {
			if idx == 0 || line[idx-1] != '\\' {
				line = line[:idx]
			}
		}
		cleaned = append(cleaned, line)
	}
	text := strings.Join(cleaned, " ")

	// Strip LaTeX commands
	text = reBeginEnd.ReplaceAllString(text, " ")
	text = reTexCommand.ReplaceAllString(text, " ")

	// Remove braces
	text = strings.ReplaceAll(text, "{", " ")
	text = strings.ReplaceAll(text, "}", " ")
	text = strings.ReplaceAll(text, "[", " ")
	text = strings.ReplaceAll(text, "]", " ")

	// Count words
	count := 0
	for _, word := range strings.Fields(text) {
		// Only count sequences containing at least one letter
		hasLetter := false
		for _, r := range word {
			if unicode.IsLetter(r) {
				hasLetter = true
				break
			}
		}
		if hasLetter {
			count++
		}
	}

	return count
}

// BibMetadata parses a .bib file and extracts structured entries.
func BibMetadata(path string) ([]BibEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening bib file: %w", err)
	}
	defer file.Close()

	var entries []BibEntry
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024) // 1MB buffer

	var currentEntry *BibEntry
	var entryContent strings.Builder
	braceDepth := 0

	for scanner.Scan() {
		line := scanner.Text()

		if currentEntry == nil {
			// Look for entry start
			if m := reBibEntry.FindStringSubmatch(line); m != nil {
				currentEntry = &BibEntry{
					Type: strings.ToLower(m[1]),
					Key:  m[2],
				}
				entryContent.Reset()
				braceDepth = 1 // We've seen the opening brace
				// Count any additional braces on this line after the match
				matchEnd := reBibEntry.FindStringIndex(line)[1]
				for _, ch := range line[matchEnd:] {
					if ch == '{' {
						braceDepth++
					} else if ch == '}' {
						braceDepth--
					}
				}
				entryContent.WriteString(line)
				entryContent.WriteString("\n")

				if braceDepth <= 0 {
					// Single-line entry
					parseEntryFields(currentEntry, entryContent.String())
					entries = append(entries, *currentEntry)
					currentEntry = nil
				}
			}
		} else {
			entryContent.WriteString(line)
			entryContent.WriteString("\n")

			for _, ch := range line {
				if ch == '{' {
					braceDepth++
				} else if ch == '}' {
					braceDepth--
				}
			}

			if braceDepth <= 0 {
				parseEntryFields(currentEntry, entryContent.String())
				entries = append(entries, *currentEntry)
				currentEntry = nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning bib file: %w", err)
	}

	return entries, nil
}

// parseEntryFields extracts title, author, year, journal from the raw entry text.
func parseEntryFields(entry *BibEntry, text string) {
	for _, m := range reBibField.FindAllStringSubmatch(text, -1) {
		field := strings.ToLower(m[1])
		value := strings.TrimSpace(m[2])
		switch field {
		case "title":
			entry.Title = value
		case "author":
			entry.Authors = value
		case "year":
			entry.Year = value
		case "journal":
			entry.Journal = value
		}
	}
}

// ImageMetadata reads image dimensions and file size.
// Supports PNG, JPEG, and GIF via Go's stdlib image decoders.
func ImageMetadata(path string) (*ImageMeta, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening image: %w", err)
	}
	defer file.Close()

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat image: %w", err)
	}

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		// If we can't decode, still return size info
		return &ImageMeta{
			Size:   stat.Size(),
			Format: "unknown",
		}, nil
	}

	return &ImageMeta{
		Width:  config.Width,
		Height: config.Height,
		Size:   stat.Size(),
		Format: format,
	}, nil
}

// FormatSize returns a human-readable file size.
func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
