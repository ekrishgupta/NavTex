package core

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// CompileResult holds the result of a LaTeX compilation.
type CompileResult struct {
	Success  bool
	LogPath  string
	Duration time.Duration
	Output   string
	Engine   string
}

// Compiler orchestrates the LaTeX build process.
// It handles multi-pass compilation and concurrency control.
type Compiler struct {
	mu   sync.Mutex
	busy bool
}

// NewCompiler creates a new thread-safe compiler manager.
func NewCompiler() *Compiler {
	return &Compiler{}
}

// IsBusy returns whether a build is currently in progress.
func (c *Compiler) IsBusy() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.busy
}

// Compile runs the LaTeX engine on the given .tex file.
// If .bib files exist in the same directory, it performs a full build:
// engine → bibtex → engine → engine
func (c *Compiler) Compile(texPath string, engine string) (*CompileResult, error) {
	c.mu.Lock()
	if c.busy {
		c.mu.Unlock()
		return nil, fmt.Errorf("a build is already in progress")
	}
	c.busy = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.busy = false
		c.mu.Unlock()
	}()

	if engine == "" {
		engine = "pdflatex"
	}

	absPath, err := filepath.Abs(texPath)
	if err != nil {
		return nil, fmt.Errorf("resolving path: %w", err)
	}

	dir := filepath.Dir(absPath)
	baseName := strings.TrimSuffix(filepath.Base(absPath), filepath.Ext(absPath))
	logPath := filepath.Join(dir, baseName+".log")

	start := time.Now()

	// Check if engine is available
	if _, err := exec.LookPath(engine); err != nil {
		return nil, fmt.Errorf("%s not found in PATH: %w", engine, err)
	}

	// First pass
	args := []string{
		"-interaction=nonstopmode",
		"-halt-on-error",
		"-file-line-error",
		absPath,
	}

	var output strings.Builder

	result := c.runEngine(engine, args, dir, &output)

	// Check if bibtex is needed (look for \bibliography or \addbibresource)
	if result {
		needsBibtex := c.checkBibtexNeeded(absPath)
		if needsBibtex {
			// Run bibtex
			c.runBibtex(baseName, dir, &output)
			// Second pass
			c.runEngine(engine, args, dir, &output)
			// Third pass for references
			result = c.runEngine(engine, args, dir, &output)
		}
	}

	duration := time.Since(start)

	return &CompileResult{
		Success:  result,
		LogPath:  logPath,
		Duration: duration,
		Output:   output.String(),
		Engine:   engine,
	}, nil
}

// runEngine executes the LaTeX engine and returns success status.
func (c *Compiler) runEngine(engine string, args []string, dir string, output *strings.Builder) bool {
	cmd := exec.Command(engine, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	output.WriteString(string(out))
	output.WriteString("\n---\n")

	return err == nil
}

// runBibtex executes bibtex on the given base name.
func (c *Compiler) runBibtex(baseName string, dir string, output *strings.Builder) {
	bibtex := "bibtex"
	if _, err := exec.LookPath(bibtex); err != nil {
		output.WriteString("bibtex not found, skipping\n")
		return
	}

	cmd := exec.Command(bibtex, baseName)
	cmd.Dir = dir

	out, _ := cmd.CombinedOutput()
	output.WriteString(string(out))
	output.WriteString("\n---\n")
}

// checkBibtexNeeded reads the .tex file to see if it uses bibliography commands.
func (c *Compiler) checkBibtexNeeded(texPath string) bool {
	data, err := readFileContent(texPath)
	if err != nil {
		return false
	}

	return strings.Contains(data, `\bibliography{`) ||
		strings.Contains(data, `\addbibresource{`) ||
		strings.Contains(data, `\printbibliography`)
}

// readFileContent reads a file and returns its content as a string.
func readFileContent(path string) (string, error) {
	data, err := exec.Command("cat", path).Output()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// OpenPDF opens the given PDF file in the system's default viewer.
func OpenPDF(pdfPath string) error {
	absPath, err := filepath.Abs(pdfPath)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	// Use platform-appropriate open command
	var cmd *exec.Cmd
	cmd = exec.Command("open", absPath) // macOS

	// Try xdg-open for Linux
	if _, err := exec.LookPath("open"); err != nil {
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", absPath)
		} else {
			return fmt.Errorf("no PDF viewer found (tried 'open' and 'xdg-open')")
		}
	}

	return cmd.Start()
}
