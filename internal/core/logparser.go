package core

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// LogEntry represents a single parsed error or warning from a TeX log.
type LogEntry struct {
	Severity string // "error" or "warning"
	Line     int    // Source line number, 0 if unknown
	Message  string
	File     string // Source file, empty if unknown
}

var (
	// TeX error: ! LaTeX Error: ...
	reTexError = regexp.MustCompile(`^!\s*(.+)`)
	// Line reference: l.123 ...
	reLineLoc = regexp.MustCompile(`^l\.(\d+)\s*(.*)`)
	// File-line-error format: ./main.tex:42: error message
	reFileLineError = regexp.MustCompile(`^(.+\.tex):(\d+):\s*(.+)`)
	// LaTeX Warning
	reLatexWarning = regexp.MustCompile(`^(?:LaTeX|Package\s+\w+)\s+Warning:\s*(.+)`)
	// Overfull/Underfull box warnings
	reBoxWarning = regexp.MustCompile(`^((?:Over|Under)full\s+\\[hv]box\s+.+)`)
	// Citation/reference warnings
	reCitationWarning = regexp.MustCompile(`Citation\s+'([^']+)'\s+on\s+page\s+\d+\s+undefined`)
	reRefWarning      = regexp.MustCompile(`Reference\s+'([^']+)'\s+on\s+page\s+\d+\s+undefined`)
	// Input file tracking
	reInputFile = regexp.MustCompile(`\(([^\s()]+\.tex)`)
)

// ParseLog reads a TeX .log file and returns structured error/warning entries.
func ParseLog(path string) ([]LogEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening log file: %w", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	var currentFile string
	var pendingError string

	for scanner.Scan() {
		line := scanner.Text()

		// Track current file
		if m := reInputFile.FindStringSubmatch(line); m != nil {
			currentFile = m[1]
		}

		// File-line-error format (triggered by -file-line-error flag)
		if m := reFileLineError.FindStringSubmatch(line); m != nil {
			lineNum, _ := strconv.Atoi(m[2])
			entries = append(entries, LogEntry{
				Severity: "error",
				Line:     lineNum,
				Message:  strings.TrimSpace(m[3]),
				File:     m[1],
			})
			pendingError = ""
			continue
		}

		// TeX error: ! ...
		if m := reTexError.FindStringSubmatch(line); m != nil {
			pendingError = strings.TrimSpace(m[1])
			continue
		}

		// Line location after an error
		if pendingError != "" {
			if m := reLineLoc.FindStringSubmatch(line); m != nil {
				lineNum, _ := strconv.Atoi(m[1])
				msg := pendingError
				if context := strings.TrimSpace(m[2]); context != "" {
					msg += " — " + context
				}
				entries = append(entries, LogEntry{
					Severity: "error",
					Line:     lineNum,
					Message:  msg,
					File:     currentFile,
				})
				pendingError = ""
				continue
			}
		}

		// If we had a pending error but no line location, still record it
		if pendingError != "" && !strings.HasPrefix(line, " ") && line != "" {
			entries = append(entries, LogEntry{
				Severity: "error",
				Line:     0,
				Message:  pendingError,
				File:     currentFile,
			})
			pendingError = ""
		}

		// LaTeX/Package warnings
		if m := reLatexWarning.FindStringSubmatch(line); m != nil {
			msg := strings.TrimSpace(m[1])
			// Multi-line warnings end with a period on a subsequent line
			msg = strings.TrimSuffix(msg, ".")
			entries = append(entries, LogEntry{
				Severity: "warning",
				Line:     extractLineFromWarning(msg),
				Message:  msg,
				File:     currentFile,
			})
			continue
		}

		// Box warnings
		if m := reBoxWarning.FindStringSubmatch(line); m != nil {
			entries = append(entries, LogEntry{
				Severity: "warning",
				Line:     0,
				Message:  strings.TrimSpace(m[1]),
				File:     currentFile,
			})
			continue
		}

		// Citation warnings
		if m := reCitationWarning.FindStringSubmatch(line); m != nil {
			entries = append(entries, LogEntry{
				Severity: "warning",
				Line:     0,
				Message:  fmt.Sprintf("Undefined citation: %s", m[1]),
				File:     currentFile,
			})
			continue
		}

		// Reference warnings
		if m := reRefWarning.FindStringSubmatch(line); m != nil {
			entries = append(entries, LogEntry{
				Severity: "warning",
				Line:     0,
				Message:  fmt.Sprintf("Undefined reference: %s", m[1]),
				File:     currentFile,
			})
			continue
		}
	}

	// Flush any remaining pending error
	if pendingError != "" {
		entries = append(entries, LogEntry{
			Severity: "error",
			Line:     0,
			Message:  pendingError,
			File:     currentFile,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning log file: %w", err)
	}

	return entries, nil
}

// ErrorCount returns the number of errors in a log entry list.
func ErrorCount(entries []LogEntry) int {
	count := 0
	for _, e := range entries {
		if e.Severity == "error" {
			count++
		}
	}
	return count
}

// WarningCount returns the number of warnings in a log entry list.
func WarningCount(entries []LogEntry) int {
	count := 0
	for _, e := range entries {
		if e.Severity == "warning" {
			count++
		}
	}
	return count
}

// extractLineFromWarning tries to find "on input line N" in a warning message.
var reWarningLine = regexp.MustCompile(`on input line\s+(\d+)`)

func extractLineFromWarning(msg string) int {
	if m := reWarningLine.FindStringSubmatch(msg); m != nil {
		n, _ := strconv.Atoi(m[1])
		return n
	}
	return 0
}
