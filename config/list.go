package config

import (
	"os"
	"strings"
)

// ListProjects returns a list of all project names
func ListProjects() ([]string, error) {
	dir := ConfigDir()
	
	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return nil, err
	}
	
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	
	var projects []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			// Remove .yaml extension to get project name
			projectName := strings.TrimSuffix(file.Name(), ".yaml")
			projects = append(projects, projectName)
		}
	}
	
	return projects, nil
}