package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/darinhaener/zellijinator/config"
	"github.com/darinhaener/zellijinator/internal/styles"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	forceDelete bool
	killSession bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete [project]",
	Short: "Delete a zellijinator project",
	Long:  `Delete the specified zellijinator project configuration`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string
		
		if len(args) == 0 {
			// No project specified, show interactive selection
			selected, err := selectProject("Select a project to delete:")
			if err != nil {
				fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("%v", err)))
				os.Exit(1)
			}
			projectName = selected
		} else {
			projectName = args[0]
		}
		
		deleteProject(projectName)
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Force deletion without confirmation")
	deleteCmd.Flags().BoolVarP(&killSession, "kill", "k", false, "Kill associated Zellij session if running")
	rootCmd.AddCommand(deleteCmd)
}

func deleteProject(name string) {
	// Get project path
	projectPath := config.ProjectPath(name)
	
	// Check if project exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Project '%s' not found.", name)))
		os.Exit(1)
	}
	
	// Read project to get session name
	data, err := os.ReadFile(projectPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error reading project file: %v", err)))
		os.Exit(1)
	}
	
	var project config.Project
	if err := yaml.Unmarshal(data, &project); err != nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error parsing project file: %v", err)))
		os.Exit(1)
	}
	
	sessionName := project.SessionName
	if sessionName == "" {
		sessionName = project.Name
	}
	
	// Check if session is running
	sessionRunning := false
	checkCmd := exec.Command("zellij", "list-sessions", "-n")
	output, err := checkCmd.CombinedOutput()
	if err == nil && len(output) > 0 {
		sessions := string(output)
		for _, line := range strings.Split(sessions, "\n") {
			if strings.TrimSpace(line) == sessionName {
				sessionRunning = true
				break
			}
		}
	}
	
	// Confirm deletion if not forced
	if !forceDelete {
		fmt.Print(styles.Prompt.Render(fmt.Sprintf("Are you sure you want to delete project '%s'? (y/N): ", name)))
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		
		if response != "y" && response != "yes" {
			fmt.Println(styles.InfoMsg("Deletion cancelled."))
			return
		}
	}
	
	// Kill session if requested and running
	if killSession && sessionRunning {
		fmt.Println(styles.InfoMsg(fmt.Sprintf("Killing Zellij session '%s'...", sessionName)))
		killCmd := exec.Command("zellij", "kill-session", sessionName)
		if err := killCmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, styles.WarningMsg(fmt.Sprintf("Failed to kill session: %v", err)))
		}
	} else if sessionRunning && !killSession {
		fmt.Println(styles.WarningMsg(fmt.Sprintf("Zellij session '%s' is still running. Use -k flag to kill it.", sessionName)))
	}
	
	// Delete the project file
	if err := os.Remove(projectPath); err != nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error deleting project file: %v", err)))
		os.Exit(1)
	}
	
	fmt.Println(styles.SuccessMsg(fmt.Sprintf("Project '%s' deleted successfully.", name)))
	
	// Clean up temporary layout files
	tmpDir := os.TempDir()
	layoutDir := strings.Join([]string{tmpDir, "zellijinator"}, string(os.PathSeparator))
	cleanupOldLayouts(layoutDir)
}