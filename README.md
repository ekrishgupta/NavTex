# NavTex TUI

A terminal-based LaTeX workspace manager built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

NavTex replaces noisy flat file listings with a **context-aware, three-pane interface** that understands `.tex` project lifecycles — hiding auxiliary clutter, surfacing metadata, and providing one-key LaTeX compilation.

## Features

| Feature | Description |
|---------|-------------|
| **Smart File Browser** | Files grouped by role: Source, Data, Assets, Shadow Bin (auxiliaries) |
| **Metadata Inspector** | Preamble & word count for `.tex`, formatted citations for `.bib`, dimensions for images |
| **One-Key Compile** | `F5` runs pdflatex/lualatex/xelatex with build status |
| **Error Log Parser** | Failed builds show clean error list (line + message) instead of raw logs |
| **One-Click Cleanup** | `F6` purges all auxiliary files instantly |
| **Template Injection** | `n` scaffolds a new project (main.tex + refs.bib + images/) |
| **Environment Sync** | Prevents concurrent builds with a mutex guard |

## Installation

```bash
go install github.com/ekrishgupta/navtex@latest
```

Or build from source:

```bash
git clone https://github.com/ekrishgupta/navtex.git
cd navtex
go build -o navtex .
```

## Usage

```bash
# Open current directory
navtex

# Open a specific project
navtex /path/to/latex/project

# Use a specific compiler
navtex --engine lualatex
```

### Keybindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move cursor up (Browser/Inspector) |
| `↓` / `j` | Move cursor down (Browser/Inspector) |
| `Tab` | Switch focus between Browser and Inspector panes |
| `F5` | Compile LaTeX document |
| `F6` | Clean auxiliary files |
| `F7` | Open compiled PDF in system viewer |
| `y` | Yank BibTeX citekey to system clipboard (when focused on Data) |
| `n` | New project wizard |
| `h` | Toggle Shadow Bin (auxiliary files) |
| `?` | Toggle Help overlay |
| `q` | Quit |

### Smart File Categories

- **Source**: `.tex` files.
- **Data**: `.bib`, `.csv`, `.py`, `.r` (for reproducible research).
- **Assets**: Images and non-LaTeX PDFs.
- **Output**: The compiled PDF matching your source.
- **Shadow Bin**: Auxiliary files (`.aux`, `.log`, `.synctex.gz`, etc.), hidden by default.

### Rich Metadata & BibMan Integration

NavTex features a rich metadata inspector:
- **BibTeX**: Displays entries in a premium tabular layout with Authors, Title, Year, and Type. Supports navigating and yanking citekeys.
- **TeX**: Extracts and tags document classes and packages.
- **Images**: Instant pixel dimensions and file size.

## Requirements

- Go 1.21+
- A LaTeX distribution (TeX Live, MiKTeX, or MacTeX) on your `$PATH`

## License

MIT
