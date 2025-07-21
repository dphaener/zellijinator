package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dphaener/zellijinator/config"
	"github.com/dphaener/zellijinator/internal/styles"
	"github.com/dphaener/zellijinator/internal/zellij"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var startCmd = &cobra.Command{
	Use:   "start [project]",
	Short: "Start a zellijinator project",
	Long:  `Start a Zellij session using the specified project configuration`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string
		
		if len(args) == 0 {
			// No project specified, show interactive selection
			selected, err := selectProject("Select a project to start:")
			if err != nil {
				fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error: %v", err)))
				os.Exit(1)
			}
			projectName = selected
		} else {
			projectName = args[0]
		}
		
		StartProject(projectName)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

// StartProject starts a Zellij session for the given project
// Exported so it can be used by root command
func StartProject(name string) {
	// Load project configuration
	projectPath := config.ProjectPath(name)
	data, err := os.ReadFile(projectPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Project %s not found. Create it with: %s", styles.Bold.Render(name), styles.Command.Render(fmt.Sprintf("zellijinator new %s", name)))))
		} else {
			fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error reading project file: %v", err)))
		}
		os.Exit(1)
	}

	// Parse YAML
	var project config.Project
	if err := yaml.Unmarshal(data, &project); err != nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error parsing project file: %v", err)))
		os.Exit(1)
	}

	// Expand home directory in root path
	if project.Root != "" && strings.HasPrefix(project.Root, "~/") {
		home, _ := os.UserHomeDir()
		project.Root = filepath.Join(home, project.Root[2:])
	}
	
	// Expand environment variables in root path
	project.Root = os.ExpandEnv(project.Root)

	// Use project name if session name not specified
	sessionName := project.SessionName
	if sessionName == "" {
		sessionName = project.Name
	}

	// Check if we're already in a Zellij session
	if os.Getenv("ZELLIJ") != "" {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg("Already inside a Zellij session."))
		fmt.Fprintln(os.Stderr, styles.InfoMsg("Please exit the current session first (Ctrl+q) before starting a new one."))
		os.Exit(1)
	}

	// Check if session already exists (active or dead)
	sessionActive := false
	checkCmd := exec.Command("zellij", "list-sessions", "-n")
	output, err := checkCmd.CombinedOutput()
	
	// If the command succeeds and we have output, check for existing session
	if err == nil && len(output) > 0 {
		sessions := string(output)
		for _, line := range strings.Split(sessions, "\n") {
			if strings.TrimSpace(line) == sessionName {
				sessionActive = true
				break
			}
		}
	}

	// If session is active, attach to it
	if sessionActive {
		fmt.Println(styles.InfoMsg(fmt.Sprintf("Attaching to existing session %s...", styles.Bold.Render(sessionName))))
		attachCmd := exec.Command("zellij", "attach", sessionName)
		attachCmd.Stdin = os.Stdin
		attachCmd.Stdout = os.Stdout
		attachCmd.Stderr = os.Stderr
		
		if err := attachCmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error attaching to session: %v", err)))
			os.Exit(1)
		}
		return
	}
	
	// Session is not in active list. Try to create it - if it exists but is dead,
	// Zellij will tell us and we'll handle it
	fmt.Println(styles.InfoMsg(fmt.Sprintf("Creating new session %s...", styles.Bold.Render(sessionName))))

	// Generate layout or use custom layout file
	var layoutPath string
	if project.Layout != "" {
		// Expand home directory in layout path
		if strings.HasPrefix(project.Layout, "~/") {
			home, _ := os.UserHomeDir()
			project.Layout = filepath.Join(home, project.Layout[2:])
		}
		layoutPath = project.Layout
	} else {
		// Generate layout from config
		layout := zellij.GenerateLayout(&project)
		
		// Create temporary layout file in a more persistent location
		tmpDir := filepath.Join(os.TempDir(), "zellijinator")
		os.MkdirAll(tmpDir, 0755)
		
		tmpFile, err := os.CreateTemp(tmpDir, fmt.Sprintf("%s-*.kdl", sessionName))
		if err != nil {
			fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error creating temp layout: %v", err)))
			os.Exit(1)
		}
		// Don't remove the file immediately - Zellij needs it!
		// We'll clean up old files on next run
		
		if _, err := tmpFile.WriteString(layout); err != nil {
			fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error writing layout: %v", err)))
			os.Exit(1)
		}
		tmpFile.Close()
		
		layoutPath = tmpFile.Name()
		
		// Clean up old layout files (older than 24 hours)
		cleanupOldLayouts(tmpDir)
		
		// Debug: print layout if ZELLIJINATOR_DEBUG is set
		if os.Getenv("ZELLIJINATOR_DEBUG") != "" {
			fmt.Println(styles.InfoMsg(fmt.Sprintf("Generated layout file: %s", styles.Path.Render(layoutPath))))
			fmt.Println(styles.InfoMsg("Layout content:"))
			fmt.Println(layout)
		}
	}

	// For new sessions, just use the layout
	// Adding --session seems to cause issues
	args := []string{"--layout", layoutPath}
	
	// Try to start the session
	startCmd := exec.Command("zellij", args...)
	startCmd.Stdin = os.Stdin
	startCmd.Stdout = os.Stdout  
	startCmd.Stderr = os.Stderr
	startCmd.Dir = project.Root
	startCmd.Env = os.Environ()
	for k, v := range project.Env {
		startCmd.Env = append(startCmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	
	if err := startCmd.Run(); err != nil {
		// If it failed, it might be because the session already exists
		// Let's check by trying to attach
		attachCmd := exec.Command("zellij", "attach", sessionName)
		attachCmd.Stdin = os.Stdin
		attachCmd.Stdout = os.Stdout
		attachCmd.Stderr = os.Stderr
		
		attachErr := attachCmd.Run()
		if attachErr == nil {
			// Successfully attached to existing session
			return
		}
		
		// Both creating and attaching failed
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error starting Zellij session: %v", err)))
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Also failed to attach to session %s: %v", styles.Bold.Render(sessionName), attachErr)))
		fmt.Fprintln(os.Stderr, styles.InfoMsg(fmt.Sprintf("Command was: %s", styles.Command.Render(fmt.Sprintf("zellij %s", strings.Join(args, " "))))))
		fmt.Fprintln(os.Stderr, styles.InfoMsg("\nDebug info:"))
		fmt.Fprintln(os.Stderr, styles.InfoMsg(fmt.Sprintf("- Layout file: %s", styles.Path.Render(layoutPath))))
		fmt.Fprintln(os.Stderr, styles.InfoMsg(fmt.Sprintf("- Session name: %s", styles.Bold.Render(sessionName))))
		fmt.Fprintln(os.Stderr, styles.InfoMsg(fmt.Sprintf("- Working directory: %s", styles.Path.Render(project.Root))))
		os.Exit(1)
	}
}

// cleanupOldLayouts removes layout files older than 24 hours
func cleanupOldLayouts(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return // Ignore errors during cleanup
	}

	cutoff := time.Now().Add(-24 * time.Hour)
	
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".kdl") {
			continue
		}
		
		info, err := file.Info()
		if err != nil {
			continue
		}
		
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(dir, file.Name()))
		}
	}
}