package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/latex"
)

// FileSelectedMsg is emitted when the cursor moves to a new file.
type FileSelectedMsg struct {
	Path     string
	Category latex.FileCategory
}

// FileBrowser is the left-pane file browser widget.
type FileBrowser struct {
	files      *latex.ProjectFiles
	items      []browserItem
	cursor     int
	width      int
	height     int
	showShadow bool
	focused    bool
	filter     string
	style      BrowserBaseStyle
}

type browserItem struct {
	display  string
	path     string
	category latex.FileCategory
	isHeader bool
}

// NewFileBrowser creates a new file browser.
func NewFileBrowser(style BrowserBaseStyle) FileBrowser {
	return FileBrowser{
		showShadow: false,
		focused:    true,
		filter:     "",
		style:      style,
	}
}

// SetStyle switches the browser's visual style (focused / blurred).
func (fb *FileBrowser) SetStyle(s BrowserBaseStyle) {
	fb.style = s
}

// SetFilter sets the search filter string.
func (fb *FileBrowser) SetFilter(f string) {
	fb.filter = strings.ToLower(f)
	fb.rebuildItems()
	fb.cursor = 0
	fb.advanceCursorToSelectable(1)
}

// SetFiles populates the browser with scanned project files.
func (fb *FileBrowser) SetFiles(f *latex.ProjectFiles) {
	fb.files = f
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
func (fb *FileBrowser) SelectedFile() (string, latex.FileCategory) {
	if len(fb.items) == 0 {
		return "", 0
	}
	item := fb.items[fb.cursor]
	return item.path, item.category
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
	fb.items = fb.items[:0]

	if fb.files == nil {
		return
	}

	// Helper to check if a file name matches the current filter
	matchFilter := func(name string) bool {
		if fb.filter == "" {
			return true
		}
		return strings.Contains(strings.ToLower(name), fb.filter)
	}

	// Source files
	if len(fb.files.Source) > 0 {
		filtered := make([]browserItem, 0, len(fb.files.Source))
		for _, f := range fb.files.Source {
			if matchFilter(f.Name) {
				filtered = append(filtered, browserItem{
					display:  f.Name,
					path:     f.Path,
					category: latex.CategorySource,
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

	// Data (includes .bib)
	{
		var filtered []browserItem
		for _, f := range fb.files.Data {
			if matchFilter(f.Name) {
				filtered = append(filtered, browserItem{
					display:  f.Name,
					path:     f.Path,
					category: latex.CategoryData,
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
	{
		var filtered []browserItem
		for _, f := range fb.files.Assets {
			rel, _ := filepath.Rel(fb.files.Root, f.Path)
			if matchFilter(rel) {
				filtered = append(filtered, browserItem{
					display:  rel,
					path:     f.Path,
				 	category: latex.CategoryAssets,
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
		filtered := make([]browserItem, 0, len(fb.files.Output))
		for _, f := range fb.files.Output {
			if matchFilter(f.Name) {
				filtered = append(filtered, browserItem{
					display:  f.Name,
					path:     f.Path,
					category: latex.CategoryOutput,
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
		filtered := make([]browserItem, 0, len(fb.files.Auxiliary))
		if fb.showShadow {
			for _, f := range fb.files.Auxiliary {
				if matchFilter(f.Name) {
					filtered = append(filtered, browserItem{
						display:  f.Name,
						path:     f.Path,
						category: latex.CategoryAuxiliary,
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
	s := fb.style
	innerW := fb.width - 2 // margin for content

	if fb.files == nil {
		empty := lipgloss.Place(innerW, fb.height-3, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center,
				LogoStyle.Render(Logo),
				"",
				DimText.Render("No project loaded"),
				DimText.Render("Press 'n' to create one"),
			),
		)
		titleBar := s.TitleBar.Width(innerW).Render("Browser")
		return lipgloss.NewStyle().Width(fb.width).Height(fb.height).Margin(0, 1).
			Render(lipgloss.JoinVertical(lipgloss.Left, titleBar, empty))
	}

	// Title bar
	titleBar := s.TitleBar.Width(innerW).Render("📄 Browser")

	var lines []string

	// Compute visible window
	visibleHeight := fb.height - 3 // title bar + margins
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
			lines = append(lines, s.CategoryHeader.Render(item.display))
		} else if i == fb.cursor && fb.focused {
			name := truncate(item.display, innerW-6)
			lines = append(lines, s.SelectedItem.Render(s.SelectedPrefix+name))
		} else if item.category == latex.CategoryAuxiliary {
			lines = append(lines, s.DimItem.Render(truncate(item.display, innerW-6)))
		} else {
			lines = append(lines, s.UnselectedItem.Render(truncate(item.display, innerW-6)))
		}
	}

	content := strings.Join(lines, "\n")
	return lipgloss.NewStyle().Width(fb.width).Height(fb.height).Margin(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, titleBar, content))
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
