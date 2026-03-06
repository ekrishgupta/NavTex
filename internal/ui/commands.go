package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekrishgupta/navtex/internal/core"
)

// ── Commands ──

func (m Model) scanDirCmd(root string) tea.Cmd {
	return func() tea.Msg {
		pf, err := core.ScanDirectory(root)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return ScannedMsg{Files: pf}
	}
}
