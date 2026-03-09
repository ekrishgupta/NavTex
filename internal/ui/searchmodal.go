package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/core"
)

// SearchModal allows fuzzy searching through the global bibliography.
type SearchModal struct {
	visible bool
	input   textinput.Model
	entries []core.BibEntry
	results []core.BibEntry
	index   int
	status  string
}

// NewSearchModal creates a new search modal.
func NewSearchModal() SearchModal {
	ti := textinput.New()
	ti.Placeholder = "Search global bibliography..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	return SearchModal{
		input:  ti,
		status: "Ready",
	}
}

// Show opens the search modal and triggers a bib load.
func (sm *SearchModal) Show(entries []core.BibEntry) {
	sm.visible = true
	sm.entries = entries
	sm.results = entries
	sm.index = 0
	sm.input.Reset()
	sm.input.Focus()
}

// Hide closes the search modal.
func (sm *SearchModal) Hide() {
	sm.visible = false
	sm.input.Blur()
}

// IsVisible returns whether the modal is shown.
func (sm *SearchModal) IsVisible() bool {
	return sm.visible
}

// HandleKey processes key messages when the modal is active.
func (sm *SearchModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if !sm.visible {
		return nil
	}

	switch msg.Type {
	case tea.KeyEscape:
		sm.Hide()
		return nil

	case tea.KeyEnter:
		if len(sm.results) > 0 {
			key := sm.results[sm.index].Key
			cite := fmt.Sprintf("\\cite{%s}", key)
			YankToClipboard(cite)
			sm.Hide()
			return nil
		}

	case tea.KeyUp, tea.KeyCtrlK:
		if sm.index > 0 {
			sm.index--
		}

	case tea.KeyDown, tea.KeyCtrlJ:
		if sm.index < len(sm.results)-1 {
			sm.index++
		}
	}

	var cmd tea.Cmd
	sm.input, cmd = sm.input.Update(msg)
	sm.filterEntries()
	return cmd
}

func (sm *SearchModal) filterEntries() {
	query := strings.ToLower(sm.input.Value())
	if query == "" {
		sm.results = sm.entries
	} else {
		var filtered []core.BibEntry
		for _, e := range sm.entries {
			if strings.Contains(strings.ToLower(e.Key), query) ||
				strings.Contains(strings.ToLower(e.Title), query) ||
				strings.Contains(strings.ToLower(e.Authors), query) {
				filtered = append(filtered, e)
			}
		}
		sm.results = filtered
	}

	if sm.index >= len(sm.results) {
		sm.index = 0
		if len(sm.results) > 0 {
			sm.index = 0
		}
	}
}

// View renders the search modal.
func (sm SearchModal) View(termWidth, termHeight int) string {
	if !sm.visible {
		return ""
	}

	title := ModalTitle.Render("Global Bibliography Search")
	inputView := sm.input.View()

	var rows []string
	maxRows := 10
	start := 0
	if sm.index >= maxRows {
		start = sm.index - maxRows + 1
	}

	for i := start; i < len(sm.results) && i < start+maxRows; i++ {
		entry := sm.results[i]
		style := FileItem
		if i == sm.index {
			style = FileItemSelected
		}

		key := style.Render(fmt.Sprintf("%-15s", entry.Key))
		title := MetaValue.Render(truncate(entry.Title, 40))
		rows = append(rows, "  "+key+" "+title)
	}

	if len(sm.results) == 0 {
		rows = append(rows, "  "+FileItemDim.Render("No matches found"))
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		"",
		title,
		"",
		"  "+inputView,
		"",
		strings.Join(rows, "\n"),
		"",
		FileItemDim.Render(fmt.Sprintf("  %d / %d entries | Enter to copy \\cite{...}", sm.index+1, len(sm.results))),
	)

	modal := ModalBox.Width(60).Render(content)
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, modal)
}
