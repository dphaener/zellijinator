package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/darinhaener/zellijinator/config"
	"github.com/darinhaener/zellijinator/internal/styles"
)

// selectProject presents an interactive list of projects using huh
// Returns the selected project name or an error
func selectProject(prompt string) (string, error) {
	// Get list of projects
	projects, err := config.ListProjects()
	if err != nil {
		return "", fmt.Errorf("error listing projects: %v", err)
	}
	
	if len(projects) == 0 {
		return "", fmt.Errorf("no projects found. Create one with: zellijinator new <project-name>")
	}
	
	// Create options for the select
	options := make([]huh.Option[string], len(projects))
	for i, project := range projects {
		options[i] = huh.NewOption(project, project)
	}
	
	// Use huh to select a project
	var selected string
	
	// Create a custom theme that matches our styling
	theme := huh.ThemeCharm()
	theme.Focused.Base = theme.Focused.Base.BorderForeground(styles.Info.GetForeground())
	theme.Focused.Title = styles.Title
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(prompt).
				Options(options...).
				Value(&selected),
		),
	).WithTheme(theme)
	
	err = form.Run()
	if err != nil {
		// User cancelled selection (e.g., pressed Ctrl-C)
		return "", fmt.Errorf("selection cancelled")
	}
	
	if selected == "" {
		return "", fmt.Errorf("no project selected")
	}
	
	return selected, nil
}