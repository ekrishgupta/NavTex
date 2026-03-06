package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/core"
)

// Inspector is the right-pane metadata viewer.
type Inspector struct {
	path     string
	category core.FileCategory
	width    int
	height   int
	focused  bool

	// Bib selection
	selectedBibIdx int

	// Cached metadata
	texMeta   *core.TexMeta
	bibMeta   []core.BibEntry
	imageMeta *core.ImageMeta
	fileSize  int64
	err       error
}

// NewInspector creates a new inspector.
func NewInspector() Inspector {
	return Inspector{}
}

// SetSize sets the inspector dimensions.
func (ins *Inspector) SetSize(w, h int) {
	ins.width = w
	ins.height = h
}

// SetFocused sets focus state.
func (ins *Inspector) SetFocused(f bool) {
	ins.focused = f
}

// SetFile updates the inspector to show metadata for the given file.
func (ins *Inspector) SetFile(path string, cat core.FileCategory) {
	if path == ins.path {
		return // No change
	}

	ins.path = path
	ins.category = cat
	ins.texMeta = nil
	ins.bibMeta = nil
	ins.imageMeta = nil
	ins.err = nil

	if path == "" {
		return
	}

	ext := strings.ToLower(filepath.Ext(path))

	switch {
	case ext == ".tex":
		meta, err := core.TexMetadata(path)
		if err != nil {
			ins.err = err
		} else {
			ins.texMeta = meta
		}

	case ext == ".bib":
		entries, err := core.BibMetadata(path)
		if err != nil {
			ins.err = err
		} else {
			ins.bibMeta = entries
		}

	case ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif":
		meta, err := core.ImageMetadata(path)
		if err != nil {
			ins.err = err
		} else {
			ins.imageMeta = meta
		}

	default:
		// For other files, just show basic info
	}
}

// MoveBibUp moves the bibliography selection up.
func (ins *Inspector) MoveBibUp() {
	if len(ins.bibMeta) == 0 {
		return
	}
	ins.selectedBibIdx--
	if ins.selectedBibIdx < 0 {
		ins.selectedBibIdx = len(ins.bibMeta) - 1
	}
}

// MoveBibDown moves the bibliography selection down.
func (ins *Inspector) MoveBibDown() {
	if len(ins.bibMeta) == 0 {
		return
	}
	ins.selectedBibIdx++
	if ins.selectedBibIdx >= len(ins.bibMeta) {
		ins.selectedBibIdx = 0
	}
}

// View renders the inspector.
func (ins Inspector) View() string {
	var content string

	if ins.path == "" {
		content = lipgloss.Place(ins.width-4, ins.height-2, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center,
				FileItemDim.Render("Select a file to inspect"),
			),
		)
	} else if ins.err != nil {
		content = lipgloss.JoinVertical(lipgloss.Left,
			InspectorTitle.Render(filepath.Base(ins.path)),
			"",
			ErrorText.Render("Error: "+ins.err.Error()),
		)
	} else {
		ext := strings.ToLower(filepath.Ext(ins.path))
		switch {
		case ext == ".tex" && ins.texMeta != nil:
			content = ins.renderTexMeta()
		case ext == ".bib" && ins.bibMeta != nil:
			content = ins.renderBibMeta()
		case (ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif") && ins.imageMeta != nil:
			content = ins.renderImageMeta()
		default:
			content = ins.renderGeneric()
		}
	}

	border := PaneBorder
	if ins.focused {
		border = PaneBorderActive
	}

	return border.Width(ins.width).Height(ins.height).Render(content)
}

// renderTexMeta renders .tex file metadata.
func (ins Inspector) renderTexMeta() string {
	m := ins.texMeta
	lines := []string{
		InspectorTitle.Render("📄 " + filepath.Base(ins.path)),
		"",
	}

	// Document info
	if m.Title != "" {
		lines = append(lines, MetaLabel.Render("Title")+" "+MetaValue.Render(m.Title))
	}
	if m.Author != "" {
		lines = append(lines, MetaLabel.Render("Author")+" "+MetaValue.Render(m.Author))
	}
	lines = append(lines, MetaLabel.Render("Class")+" "+MetaValue.Render(m.DocumentClass))
	if m.ClassOptions != "" {
		lines = append(lines, MetaLabel.Render("Options")+" "+MetaValue.Render(m.ClassOptions))
	}
	lines = append(lines, MetaLabel.Render("Word Count")+" "+MetaValue.Render(fmt.Sprintf("%d", m.WordCount)))

	// Packages
	if len(m.Packages) > 0 {
		lines = append(lines, "", CategoryLabel.Render(fmt.Sprintf("Packages (%d)", len(m.Packages))))
		var pkgLine strings.Builder
		for i, pkg := range m.Packages {
			if i > 0 {
				pkgLine.WriteString(" ")
			}
			pkgLine.WriteString(PackageTag.Render(pkg.Name))
			// Wrap long lines
			if pkgLine.Len() > ins.width-8 {
				lines = append(lines, "   "+pkgLine.String())
				pkgLine.Reset()
			}
		}
		if pkgLine.Len() > 0 {
			lines = append(lines, "   "+pkgLine.String())
		}
	}

	return strings.Join(lines, "\n")
}

// renderBibMeta renders .bib file metadata in a BibMan-inspired tabular format.
func (ins Inspector) renderBibMeta() string {
	lines := []string{
		InspectorTitle.Render("📚 " + filepath.Base(ins.path)),
		MetaValue.Render(fmt.Sprintf("   %d entries", len(ins.bibMeta))),
		"",
	}

	if len(ins.bibMeta) == 0 {
		lines = append(lines, FileItemDim.Render("  No entries found"))
		return strings.Join(lines, "\n")
	}

	// Column widths
	maxAuth := 18
	maxTitle := ins.width - maxAuth - 14 // year(4) + type(8) + padding
	if maxTitle < 20 {
		maxTitle = 20
	}

	// Header
	header := fmt.Sprintf("  %-*s %-*s %-4s %-8s",
		maxAuth, "Authors", maxTitle, "Title", "Year", "Type")
	lines = append(lines, BibTableHeader.Render(header))
	lines = append(lines, FileItemDim.Render("  "+strings.Repeat("─", min(ins.width-6, len(header)))))

	// Entries
	for i, entry := range ins.bibMeta {
		authors := truncate(entry.Authors, maxAuth)
		title := truncate(entry.Title, maxTitle)
		row := fmt.Sprintf("  %-*s %-*s %-4s %-8s",
			maxAuth, authors, maxTitle, title, entry.Year, entry.Type)

		if i == ins.selectedBibIdx && ins.focused {
			lines = append(lines, BibTableRowSelected.Width(ins.width-4).Render(row))
		} else {
			lines = append(lines, BibTableRow.Render(row))
		}

		// Show DOI/URL if present
		if entry.DOI != "" {
			lines = append(lines, FileItemDim.Render(fmt.Sprintf("    DOI: %s", entry.DOI)))
		}

		// Show keywords if present
		if len(entry.Keywords) > 0 {
			var kwLine strings.Builder
			kwLine.WriteString("    ")
			for i, kw := range entry.Keywords {
				if i > 0 {
					kwLine.WriteString(" ")
				}
				kwLine.WriteString(KeywordTag.Render(kw))
			}
			lines = append(lines, kwLine.String())
		}
	}

	return strings.Join(lines, "\n")
}

// renderImageMeta renders image file metadata.
func (ins Inspector) renderImageMeta() string {
	m := ins.imageMeta
	lines := []string{
		InspectorTitle.Render("🖼  " + filepath.Base(ins.path)),
		"",
		MetaLabel.Render("Format") + " " + MetaValue.Render(strings.ToUpper(m.Format)),
	}

	if m.Width > 0 && m.Height > 0 {
		lines = append(lines, MetaLabel.Render("Dimensions")+" "+MetaValue.Render(fmt.Sprintf("%d × %d px", m.Width, m.Height)))
	}

	lines = append(lines, MetaLabel.Render("File Size")+" "+MetaValue.Render(core.FormatSize(m.Size)))

	return strings.Join(lines, "\n")
}

// renderGeneric renders basic file information.
func (ins Inspector) renderGeneric() string {
	name := filepath.Base(ins.path)
	ext := filepath.Ext(ins.path)

	return strings.Join([]string{
		InspectorTitle.Render("📎 " + name),
		"",
		MetaLabel.Render("Extension") + " " + MetaValue.Render(ext),
		MetaLabel.Render("Category") + " " + MetaValue.Render(categoryName(ins.category)),
	}, "\n")
}

func categoryName(c core.FileCategory) string {
	switch c {
	case core.CategorySource:
		return "Source"
	case core.CategoryData:
		return "Data"
	case core.CategoryAssets:
		return "Asset"
	case core.CategoryAuxiliary:
		return "Auxiliary"
	case core.CategoryOutput:
		return "Output"
	default:
		return "Unknown"
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
