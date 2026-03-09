package core

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ProjectConfig holds settings defined in .navtex.yaml
type ProjectConfig struct {
	Engine     string   `yaml:"engine"`
	MasterFile string   `yaml:"master"`
	Ignores    []string `yaml:"ignores"`
}

// GlobalConfig holds settings defined in ~/.navtex.yaml
type GlobalConfig struct {
	GlobalBibPath string `yaml:"global_bib"`
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() ProjectConfig {
	return ProjectConfig{
		Engine:     "pdflatex",
		MasterFile: "",
	}
}

// DefaultGlobalConfig returns a global configuration with sensible defaults.
func DefaultGlobalConfig() GlobalConfig {
	return GlobalConfig{
		GlobalBibPath: "",
	}
}

// LoadConfig attempts to read .navtex.yaml from the project root.
// Returns DefaultConfig if not found.
func LoadConfig(dir string) ProjectConfig {
	config := DefaultConfig()
	path := filepath.Join(dir, ".navtex.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return config // Return default silently if missing
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config // Assume default if broken
	}

	return config
}

// LoadGlobalConfig attempts to read ~/.navtex.yaml.
func LoadGlobalConfig() GlobalConfig {
	config := DefaultGlobalConfig()
	home, err := os.UserHomeDir()
	if err != nil {
		return config
	}

	path := filepath.Join(home, ".navtex.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return config
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config
	}

	return config
}
