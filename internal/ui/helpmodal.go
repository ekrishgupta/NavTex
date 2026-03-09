package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HelpModal displays the keybinding reference.
type HelpModal struct {
	visible bool
}

// NewHelpModal creates a new help modal.
func NewHelpModal() HelpModal {
	return HelpModal{}
}

// Show opens the help modal.
func (hm *HelpModal) Show() {
	hm.visible = true
}

// Hide closes the help modal.
func (hm *HelpModal) Hide() {
	hm.visible = false
}

// IsVisible returns whether the modal is shown.
func (hm *HelpModal) IsVisible() bool {
	return hm.visible
}

// Toggle toggles the help modal.
func (hm *HelpModal) Toggle() {
	hm.visible = !hm.visible
}

// View renders the help modal.
func (hm HelpModal) View(termWidth, termHeight int) string {
	if !hm.visible {
		return ""
	}

	bindings := []struct {
		key  string
		desc string
	}{
		{"↑ / k", "Move cursor up"},
		{"↓ / j", "Move cursor down"},
		{"Tab", "Switch focus between panes"},
		{"h", "Toggle Shadow Bin (auxiliary files)"},
		{"F5", "Compile LaTeX document"},
		{"F6", "Clean auxiliary files"},
		{"F7", "Open compiled PDF"},
		{"d", "Generate latexdiff PDF"},
		{"n", "New project wizard"},

		{"y", "Yank cite key to clipboard"},
		{"s", "Search global bibliography"},
		{"/", "Search files"},
		{"Enter", "Open editor / Jump to line"},
		{"?", "Toggle this help"},
		{"q / Ctrl+C", "Quit"},
	}

	title := ModalTitle.Render("Keybindings")

	var rows []string
	for _, b := range bindings {
		key := ActionKey.Render(padRight(b.key, 12))
		desc := MetaValue.Render(b.desc)
		rows = append(rows, "  "+key+" "+desc)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		"",
		LogoStyle.Render(Logo),
		"",
		title,
		"",
		strings.Join(rows, "\n"),
		"",
		FileItemDim.Render("  Press Esc or ? to close"),
	)

	modalW := 48
	if modalW > termWidth-4 {
		modalW = termWidth - 4
	}

	modal := ModalBox.Width(modalW).Render(content)
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, modal)
}

// padRight pads a string with spaces to the given width using formatting.
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
