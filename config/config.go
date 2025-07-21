package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Project struct {
	Name          string            `yaml:"name"`
	Root          string            `yaml:"root"`
	SessionName   string            `yaml:"session_name,omitempty"`
	Layout        string            `yaml:"layout,omitempty"`
	DefaultLayout string            `yaml:"default_layout,omitempty"`
	Tabs          []Tab             `yaml:"tabs"`
	Env           map[string]string `yaml:"env,omitempty"`
}

type Tab struct {
	Name   string   `yaml:"name"`
	Focus  bool     `yaml:"focus,omitempty"`
	Layout string   `yaml:"layout,omitempty"`
	Panes  []Pane   `yaml:"panes"`
}

type Pane struct {
	Focus    bool     `yaml:"focus,omitempty"`
	Commands []string `yaml:"commands,omitempty"`
	Size     string   `yaml:"size,omitempty"`
	Split    string   `yaml:"split,omitempty"`
}

func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".zellijinator")
}

func ProjectPath(name string) string {
	return filepath.Join(ConfigDir(), fmt.Sprintf("%s.yaml", name))
}

func EnsureConfigDir() error {
	dir := ConfigDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}