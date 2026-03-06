package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekrishgupta/navtex/internal/ui"
)

var version = "0.1.0"

func main() {
	// CLI Flags
	engine := flag.String("engine", "pdflatex", "LaTeX engine to use (pdflatex, lualatex, xelatex)")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "NavTex TUI — LaTeX Workspace Manager\n\n")
		fmt.Fprintf(os.Stderr, "Usage: navtex [options] [path]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nKeybindings:\n")
		fmt.Fprintf(os.Stderr, "  F5: Compile    F6: Clean    F7: Open PDF\n")
		fmt.Fprintf(os.Stderr, "  Tab: Focus     h: Shadow    n: New Project\n")
		fmt.Fprintf(os.Stderr, "  ?: Help        q: Quit\n")
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("NavTex version %s\n", version)
		os.Exit(0)
	}

	// Positional path argument
	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	// Initialize the model
	m := ui.NewModel(path, *engine)

	// Run the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running NavTex: %v\n", err)
		os.Exit(1)
	}
}
