package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dphaener/zellijinator/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	showActive bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all zellijinator projects",
	Long:  `List all available zellijinator project configurations`,
	Run: func(cmd *cobra.Command, args []string) {
		listProjects()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&showActive, "active", "a", false, "Show active Zellij sessions")
	rootCmd.AddCommand(listCmd)
}

func listProjects() {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1)
	
	projectStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))
	
	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)
	
	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		PaddingLeft(2)
	
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))
	
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	projects, err := config.ListProjects()
	if err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render(fmt.Sprintf("Error listing projects: %v", err)))
		os.Exit(1)
	}

	if len(projects) == 0 {
		fmt.Println(dimStyle.Render("No projects found."))
		fmt.Println(dimStyle.Render("Create a new project with: zellijinator new <project-name>"))
		return
	}

	// Get active sessions
	activeSessions := make(map[string]bool)
	checkCmd := exec.Command("zellij", "list-sessions", "-n")
	output, err := checkCmd.CombinedOutput()
	if err == nil && len(output) > 0 {
		sessions := string(output)
		for _, line := range strings.Split(sessions, "\n") {
			session := strings.TrimSpace(line)
			if session != "" {
				activeSessions[session] = true
			}
		}
	}

	fmt.Println(titleStyle.Render("Zellijinator Projects"))

	for _, project := range projects {
		// Read project file to get session name and info
		projectPath := config.ProjectPath(project)
		data, err := os.ReadFile(projectPath)
		if err != nil {
			fmt.Printf("  %s %s\n", 
				projectStyle.Render(project),
				errorStyle.Render("(error reading config)"))
			continue
		}

		var proj config.Project
		if err := yaml.Unmarshal(data, &proj); err != nil {
			fmt.Printf("  %s %s\n", 
				projectStyle.Render(project),
				errorStyle.Render("(error parsing config)"))
			continue
		}

		sessionName := proj.SessionName
		if sessionName == "" {
			sessionName = proj.Name
		}

		// Format project name with status
		projectLine := "  " + projectStyle.Render(project)
		if activeSessions[sessionName] {
			projectLine += " " + activeStyle.Render("ACTIVE")
		}
		fmt.Println(projectLine)
		
		// Show additional info
		if proj.Root != "" {
			fmt.Println(infoStyle.Render(fmt.Sprintf("Root: %s", proj.Root)))
		}
		if len(proj.Tabs) > 0 {
			tabNames := make([]string, 0, len(proj.Tabs))
			for _, tab := range proj.Tabs {
				tabNames = append(tabNames, tab.Name)
			}
			fmt.Println(infoStyle.Render(fmt.Sprintf("Tabs: %s", strings.Join(tabNames, ", "))))
		}
		fmt.Println() // Add spacing between projects
	}
}