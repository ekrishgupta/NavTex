package core

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetGitLastCommitContent returns the content of a file at the last commit (HEAD).
func GetGitLastCommitContent(filePath string) (string, error) {
	return GetGitVersionContent(filePath, "HEAD")
}

// GetGitVersionContent returns the content of a file at a specific Git revision.
func GetGitVersionContent(filePath string, revision string) (string, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", revision, filePath))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git show: %w (output: %s)", err, string(out))
	}
	return string(out), nil
}

// ListGitTags returns a list of Git tags in the repository.
func ListGitTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "-l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git tag: %w (output: %s)", err, string(out))
	}
	tags := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(tags) == 1 && tags[0] == "" {
		return []string{}, nil
	}
	return tags, nil
}

// IsGitRepo checks if the current directory is a Git repository.
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}
