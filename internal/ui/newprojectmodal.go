package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ekrishgupta/navtex/internal/core"
)

// NewProjectModal is a form for scaffolding a new LaTeX project.
type NewProjectModal struct {
	visible  bool
	fields   [4]string // title, author, docclass path
	labels   [4]string
	cursor   int
	width    int
	height   int
	root     string
	inputBuf string
}

// NewNewProjectModal creates a new project modal.
func NewNewProjectModal() NewProjectModal {
	return NewProjectModal{
		labels: [4]string{"Title", "Author", "Class", "Path"},
		fields: [4]string{"", "", "article", "."},
	}
}

// Show opens the modal.
func (npm *NewProjectModal) Show(root string) {
	npm.visible = true
	npm.cursor = 0
	npm.root = root
	npm.fields[3] = root
}

// Hide closes the modal.
func (npm *NewProjectModal) Hide() {
	npm.visible = false
}

// IsVisible returns whether the modal is shown.
func (npm *NewProjectModal) IsVisible() bool {
	return npm.visible
}

// HandleKey processes a key event and returns a command if the project should be created.
func (npm *NewProjectModal) HandleKey(key tea.KeyMsg) tea.Cmd {
	switch key.Type {
	case tea.KeyEscape:
		npm.Hide()
		return nil
	case tea.KeyTab, tea.KeyDown:
		npm.cursor = (npm.cursor + 1) % 4
	case tea.KeyShiftTab, tea.KeyUp:
		npm.cursor = (npm.cursor + 3) % 4
	case tea.KeyBackspace:
		if len(npm.fields[npm.cursor]) > 0 {
			npm.fields[npm.cursor] = npm.fields[npm.cursor][:len(npm.fields[npm.cursor])-1]
		}
	case tea.KeyEnter:
		if npm.cursor == 3 {
			// Submit
			return npm.submit()
		}
		npm.cursor = (npm.cursor + 1) % 4
	case tea.KeyRunes:
		npm.fields[npm.cursor] += string(key.Runes)
	}
	return nil
}

// submit creates the project and returns a command.
func (npm *NewProjectModal) submit() tea.Cmd {
	title := npm.fields[0]
	author := npm.fields[1]
	docclass := npm.fields[2]
	path := npm.fields[3]

	if title == "" {
		title = "Untitled"
	}
	if docclass == "" {
		docclass = "article"
	}

	return func() tea.Msg {
		err := core.CreateProject(path, title, author, docclass)
		if err != nil {
			return ProjectCreatedMsg{Err: err}
		}
		return ProjectCreatedMsg{Path: path}
	}
}

// ProjectCreatedMsg is sent when a new project is created.
type ProjectCreatedMsg struct {
	Path string
	Err  error
}

// View renders the new project modal.
func (npm NewProjectModal) View(termWidth, termHeight int) string {
	if !npm.visible {
		return ""
	}

	modalW := 50
	if modalW > termWidth-4 {
		modalW = termWidth - 4
	}

	title := ModalTitle.Render("New LaTeX Project")

	var rows []string
	for i, label := range npm.labels {
		value := npm.fields[i]
		cursor := " "
		if i == npm.cursor {
			cursor = "▸"
			value += "█"
		}

		labelStr := InputLabel.Render(label + ":")
		fieldStyle := InputField
		if i == npm.cursor {
			fieldStyle = fieldStyle.BorderForeground(ColorAccent)
		}
		fieldStr := fieldStyle.Width(modalW - 16).Render(value)

		rows = append(rows, cursor+" "+labelStr+" "+fieldStr)
	}

	hint := FileItemDim.Render("  Tab: next field │ Enter on Path: create │ Esc: cancel")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		strings.Join(rows, "\n"),
		"",
		hint,
	)

	modal := ModalBox.Width(modalW).Render(content)
	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, modal)
}
