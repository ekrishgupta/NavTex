package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/core"
)

// FileSelectedMsg is emitted when the cursor moves to a new file.
type FileSelectedMsg struct {
	Path     string
	Category core.FileCategory
}

// FileBrowser is the left-pane file browser widget.
type FileBrowser struct {
	files      *core.ProjectFiles
	items      []browserItem
	cursor     int
	width      int
	height     int
	showShadow bool
	focused    bool
	filter     string
}

type browserItem struct {
	display  string
	path     string
	category core.FileCategory
	isHeader bool
}

// NewFileBrowser creates a new file browser.
func NewFileBrowser() FileBrowser {
	return FileBrowser{
		showShadow: false,
		focused:    true,
		filter:     "",
	}
}

// SetFilter sets the search filter string.
func (fb *FileBrowser) SetFilter(f string) {
	fb.filter = strings.ToLower(f)
	fb.rebuildItems()
	fb.cursor = 0
	fb.advanceCursorToSelectable(1)
}

// SetFiles populates the browser with scanned project files.
func (fb *FileBrowser) SetFiles(pf *core.ProjectFiles) {
	fb.files = pf
	fb.rebuildItems()
	if fb.cursor >= len(fb.items) {
		fb.cursor = 0
	}
	// Advance cursor to first selectable item
	fb.advanceCursorToSelectable(1)
}

// SetSize sets the dimensions of the browser.
func (fb *FileBrowser) SetSize(w, h int) {
	fb.width = w
	fb.height = h
}

// SetFocused sets focus state.
func (fb *FileBrowser) SetFocused(f bool) {
	fb.focused = f
}

// SelectedFile returns the currently selected file path and its category.
func (fb *FileBrowser) SelectedFile() (string, core.FileCategory) {
	if fb.cursor >= 0 && fb.cursor < len(fb.items) {
		item := fb.items[fb.cursor]
		if !item.isHeader {
			return item.path, item.category
		}
	}
	return "", core.CategorySource
}

// MoveUp moves cursor up.
func (fb *FileBrowser) MoveUp() {
	if len(fb.items) == 0 {
		return
	}
	fb.cursor--
	if fb.cursor < 0 {
		fb.cursor = len(fb.items) - 1
	}
	fb.advanceCursorToSelectable(-1)
}

// MoveDown moves cursor down.
func (fb *FileBrowser) MoveDown() {
	if len(fb.items) == 0 {
		return
	}
	fb.cursor++
	if fb.cursor >= len(fb.items) {
		fb.cursor = 0
	}
	fb.advanceCursorToSelectable(1)
}

// ToggleShadow toggles the shadow bin (auxiliary files) visibility.
func (fb *FileBrowser) ToggleShadow() {
	fb.showShadow = !fb.showShadow
	fb.rebuildItems()
}

// ShowingShadow returns whether the shadow bin is visible.
func (fb *FileBrowser) ShowingShadow() bool {
	return fb.showShadow
}

// rebuildItems reconstructs the flat list of display items.
func (fb *FileBrowser) rebuildItems() {
	fb.items = nil

	if fb.files == nil {
		return
	}

	// Helper to check if a file name matches the current filter
	matchFilter := func(name string) bool {
		if fb.filter == "" {
			return true
		}
		// Basic "fuzzy" match by converting to lower and checking contain
		// A full fuzzy matching algorithm could be added here later.
		return strings.Contains(strings.ToLower(name), fb.filter)
	}

	// Source files
	if len(fb.files.Source) > 0 {
		var filtered []browserItem
		for _, f := range fb.files.Source {
			if matchFilter(f.Name) {
				filtered = append(filtered, browserItem{
					display:  f.Name,
					path:     f.Path,
					category: core.CategorySource,
				})
			}
		}
		if len(filtered) > 0 {
			fb.items = append(fb.items, browserItem{
				display:  fmt.Sprintf("📄 Source (%d)", len(filtered)),
				isHeader: true,
			})
			fb.items = append(fb.items, filtered...)
		}
	}

	// Data files
	if len(fb.files.Data) > 0 {
		var filtered []browserItem
		for _, f := range fb.files.Data {
			if matchFilter(f.Name) {
				filtered = append(filtered, browserItem{
					display:  f.Name,
					path:     f.Path,
					category: core.CategoryData,
				})
			}
		}
		if len(filtered) > 0 {
			fb.items = append(fb.items, browserItem{
				display:  fmt.Sprintf("📊 Data (%d)", len(filtered)),
				isHeader: true,
			})
			fb.items = append(fb.items, filtered...)
		}
	}

	// Assets
	if len(fb.files.Assets) > 0 {
		var filtered []browserItem
		for _, f := range fb.files.Assets {
			rel, _ := filepath.Rel(fb.files.Root, f.Path)
			if matchFilter(rel) {
				filtered = append(filtered, browserItem{
					display:  rel,
					path:     f.Path,
					category: core.CategoryAssets,
				})
			}
		}
		if len(filtered) > 0 {
			fb.items = append(fb.items, browserItem{
				display:  fmt.Sprintf("🖼  Assets (%d)", len(filtered)),
				isHeader: true,
			})
			fb.items = append(fb.items, filtered...)
		}
	}

	// Output
	if len(fb.files.Output) > 0 {
		var filtered []browserItem
		for _, f := range fb.files.Output {
			if matchFilter(f.Name) {
				filtered = append(filtered, browserItem{
					display:  f.Name,
					path:     f.Path,
					category: core.CategoryOutput,
				})
			}
		}
		if len(filtered) > 0 {
			fb.items = append(fb.items, browserItem{
				display:  fmt.Sprintf("📦 Output (%d)", len(filtered)),
				isHeader: true,
			})
			fb.items = append(fb.items, filtered...)
		}
	}

	// Shadow Bin (auxiliary)
	if len(fb.files.Auxiliary) > 0 {
		var filtered []browserItem
		if fb.showShadow {
			for _, f := range fb.files.Auxiliary {
				if matchFilter(f.Name) {
					filtered = append(filtered, browserItem{
						display:  f.Name,
						path:     f.Path,
						category: core.CategoryAuxiliary,
					})
				}
			}
		}

		if len(filtered) > 0 || (!fb.showShadow && matchFilter("shadow bin")) {
			shadowLabel := fmt.Sprintf("👻 Shadow Bin (%d)", len(fb.files.Auxiliary))
			if !fb.showShadow {
				shadowLabel += " [h to show]"
			}
			fb.items = append(fb.items, browserItem{
				display:  shadowLabel,
				isHeader: true,
			})
			fb.items = append(fb.items, filtered...)
		}
	}
}

// advanceCursorToSelectable skips headers.
func (fb *FileBrowser) advanceCursorToSelectable(direction int) {
	if len(fb.items) == 0 {
		return
	}
	maxAttempts := len(fb.items)
	for i := 0; i < maxAttempts; i++ {
		if fb.cursor >= 0 && fb.cursor < len(fb.items) && !fb.items[fb.cursor].isHeader {
			return
		}
		fb.cursor += direction
		if fb.cursor < 0 {
			fb.cursor = len(fb.items) - 1
		} else if fb.cursor >= len(fb.items) {
			fb.cursor = 0
		}
	}
}

// View renders the file browser.
func (fb FileBrowser) View() string {
	if fb.files == nil {
		return PaneBorder.Width(fb.width).Height(fb.height).Render(
			lipgloss.Place(fb.width, fb.height, lipgloss.Center, lipgloss.Center,
				lipgloss.JoinVertical(lipgloss.Center,
					LogoStyle.Render(Logo),
					"",
					FileItemDim.Render("No project loaded"),
					FileItemDim.Render("Press 'n' to create one"),
				),
			),
		)
	}

	var lines []string

	// Compute visible window
	visibleHeight := fb.height - 2 // account for border padding
	scrollOffset := 0
	if fb.cursor >= visibleHeight {
		scrollOffset = fb.cursor - visibleHeight + 1
	}

	for i, item := range fb.items {
		if i < scrollOffset {
			continue
		}
		if len(lines) >= visibleHeight {
			break
		}

		if item.isHeader {
			lines = append(lines, CategoryLabel.Render(item.display))
		} else if i == fb.cursor && fb.focused {
			name := truncate(item.display, fb.width-6)
			lines = append(lines, FileItemSelected.Width(fb.width-4).Render(" "+name))
		} else if item.category == core.CategoryAuxiliary {
			lines = append(lines, FileItemDim.Render(truncate(item.display, fb.width-6)))
		} else {
			lines = append(lines, FileItem.Render(truncate(item.display, fb.width-6)))
		}
	}

	content := strings.Join(lines, "\n")

	border := PaneBorder
	if fb.focused {
		border = PaneBorderActive
	}

	return border.Width(fb.width).Height(fb.height).Render(content)
}

// truncate shortens a string to fit a given width.
func truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}
