package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/core"
)

// ErrorModal displays parsed build errors.
type ErrorModal struct {
	entries []core.LogEntry
	visible bool
	width   int
	height  int
	scroll  int
}

// NewErrorModal creates a new error modal.
func NewErrorModal() ErrorModal {
	return ErrorModal{}
}

// Show displays the modal with the given log entries.
func (em *ErrorModal) Show(entries []core.LogEntry) {
	em.entries = entries
	em.visible = true
	em.scroll = 0
}

// Hide closes the modal.
func (em *ErrorModal) Hide() {
	em.visible = false
}

// IsVisible returns whether the modal is shown.
func (em *ErrorModal) IsVisible() bool {
	return em.visible
}

// ScrollUp scrolls the error list up.
func (em *ErrorModal) ScrollUp() {
	if em.scroll > 0 {
		em.scroll--
	}
}

// ScrollDown scrolls the error list down.
func (em *ErrorModal) ScrollDown() {
	maxScroll := len(em.entries) - (em.height / 2)
	if maxScroll < 0 {
		maxScroll = 0
	}
	if em.scroll < maxScroll {
		em.scroll++
	}
}

// View renders the error modal.
func (em ErrorModal) View(termWidth, termHeight int) string {
	if !em.visible {
		return ""
	}

	modalW := termWidth * 3 / 4
	modalH := termHeight * 3 / 4
	if modalW < 60 {
		modalW = 60
	}
	if modalH < 10 {
		modalH = 10
	}

	errors := core.ErrorCount(em.entries)
	warnings := core.WarningCount(em.entries)

	title := ModalTitle.Render(fmt.Sprintf("Build Log — %d errors, %d warnings", errors, warnings))

	// Header row
	lineCol := 6
	sevCol := 8
	msgCol := modalW - lineCol - sevCol - 12

	header := fmt.Sprintf("  %-*s %-*s %s", lineCol, "Line", sevCol, "Severity", "Message")
	headerLine := BibTableHeader.Render(header)
	separator := FileItemDim.Render("  " + strings.Repeat("─", modalW-8))

	var rows []string
	for i := em.scroll; i < len(em.entries) && len(rows) < modalH-6; i++ {
		e := em.entries[i]
		lineStr := "—"
		if e.Line > 0 {
			lineStr = fmt.Sprintf("%d", e.Line)
		}

		msg := truncate(e.Message, msgCol)

		var row string
		if e.Severity == "error" {
			sev := ErrorText.Render("error")
			row = fmt.Sprintf("  %-*s %-*s %s", lineCol, lineStr, sevCol, sev, msg)
		} else {
			sev := WarningText.Render("warning")
			row = fmt.Sprintf("  %-*s %-*s %s", lineCol, lineStr, sevCol, sev, msg)
		}
		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		headerLine,
		separator,
		strings.Join(rows, "\n"),
		"",
		FileItemDim.Render("  Press Esc to close │ ↑↓ to scroll"),
	)

	modal := ModalBox.Width(modalW).Render(content)
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, modal)
}
