package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dphaener/zellijinator/config"
	"github.com/dphaener/zellijinator/internal/styles"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project]",
	Short: "Create a new zellijinator project",
	Long:  `Create a new zellijinator project configuration file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		createNewProject(projectName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func createNewProject(name string) {
	// Ensure config directory exists
	if err := config.EnsureConfigDir(); err != nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error creating config directory: %v", err)))
		os.Exit(1)
	}

	// Check if project already exists
	projectPath := config.ProjectPath(name)
	if _, err := os.Stat(projectPath); err == nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Project '%s' already exists", name)))
		fmt.Fprintln(os.Stderr, styles.Subtle.Render(fmt.Sprintf("  Location: %s", projectPath)))
		os.Exit(1)
	}

	// Create the config file with sample content
	if err := os.WriteFile(projectPath, []byte(getSampleTemplate(name)), 0644); err != nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error creating project file: %v", err)))
		os.Exit(1)
	}

	fmt.Println(styles.SuccessMsg(fmt.Sprintf("Created new zellijinator project: %s", styles.Bold.Render(name))))
	fmt.Println(styles.InfoMsg("Config file: " + styles.Path.Render(projectPath)))

	// Open in editor
	openInEditor(projectPath)
}

func openInEditor(path string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Try common editors
		for _, e := range []string{"vim", "vi", "nano", "emacs", "code"} {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}

	if editor == "" {
		fmt.Println(styles.WarningMsg("No editor found. Please set EDITOR environment variable."))
		fmt.Println(styles.InfoMsg("You can manually edit: " + styles.Path.Render(path)))
		return
	}

	fmt.Println(styles.InfoMsg(fmt.Sprintf("Opening in %s...", styles.Command.Render(editor))))
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error opening editor: %v", err)))
		fmt.Println(styles.InfoMsg("You can manually edit: " + styles.Path.Render(path)))
	}
}

func getSampleTemplate(projectName string) string {
	return fmt.Sprintf(`# Zellijinator Project Configuration
# Project: %s
#
# This file defines how your Zellij session will be created.
# Zellij is a terminal workspace with batteries included.
#
# Format: YAML
# Location: ~/.zellijinator/%s.yaml

# Project name (required)
# This is used to identify your project
name: %s

# Root directory (required)
# The directory where your project is located
# You can use ~ for home directory or environment variables like $HOME
root: ~/projects/%s

# Session name (optional)
# The name of the Zellij session. If not specified, the project name is used
# session_name: %s-session

# Layout file (optional)
# Path to a custom Zellij layout file (.kdl format)
# If not specified, zellijinator will create a layout based on the tabs configuration below
# layout: ~/.config/zellij/layouts/custom.kdl

# Default layout template (optional)
# The Zellij layout template to inherit from (e.g., "default", "compact", "strider")
# This controls UI elements like tab bar position and help text visibility
# Common options:
#   - "default": Standard layout with tab bar at top and full status bar
#   - "compact": Minimal layout with tab bar at bottom and reduced status info
#   - "strider": Layout with file browser plugin
# If not specified, uses the standard default layout
# default_layout: compact

# Environment variables (optional)
# These will be set in all panes of this session
env:
  # NODE_ENV: development
  # DATABASE_URL: postgresql://localhost:5432/mydb

# Tabs configuration (required)
# Define the tabs and panes for your session
tabs:
  # First tab example
  - name: "main"
    # Set focus to this tab on startup (only one tab should have focus: true)
    focus: true
    panes:
      # First pane - simple example
      - commands:
          - echo "Welcome to %s!"
          - ls -la

  # Second tab example with split panes
  - name: "development"
    panes:
      # First pane (main pane)
      - focus: true  # This pane will be focused when switching to this tab
        commands:
          # Commands run sequentially when the pane starts
          - echo "Starting development server..."
          # - npm run dev
        
      # Second pane - horizontal split
      - split: horizontal  # Split direction: horizontal or vertical
        size: 30  # Size in percentage (30%% of available space)
        commands:
          - echo "Watching for file changes..."
          # - npm run watch

      # Third pane - vertical split from the second pane
      - split: vertical
        size: 50  # This creates a 50-50 vertical split of the remaining space
        commands:
          - echo "Logs will appear here..."
          # - tail -f logs/app.log

  # Third tab example - database and monitoring
  - name: "backend"
    panes:
      - commands:
          - echo "Database console"
          # - psql mydb
      
      - split: horizontal
        size: 40
        commands:
          - echo "Redis monitoring"
          # - redis-cli monitor

  # Fourth tab example - simple editing tab
  - name: "editor"
    panes:
      - commands:
          - echo "Open your editor here"
          # - vim .
          # - code .

  # Fifth tab example - using predefined layouts
  - name: "monitoring"
    # Use a predefined layout instead of manual splits
    layout: tiled  # Will arrange 3 panes in a tiled pattern
    panes:
      - commands:
          - echo "System stats"
          # - htop
      - commands:
          - echo "Logs"
          # - tail -f logs/production.log
      - commands:
          - echo "Network monitoring"
          # - sudo nethogs

# CONFIGURATION GUIDE:
# 
# 1. TABS:
#    - Each tab represents a workspace within your Zellij session
#    - Tabs appear at the top of your terminal
#    - You can switch between tabs using Alt+[number] or Alt+h/l
#
# 2. TAB LAYOUTS:
#    - You can use predefined layouts for automatic pane arrangement:
#      - 'even-horizontal': All panes split horizontally with equal height
#      - 'even-vertical': All panes split vertically with equal width  
#      - 'main-vertical': Large pane on left (70%), others stacked on right
#      - 'main-horizontal': Large pane on top (70%), others side-by-side below
#      - 'tiled': Arranges panes in a grid pattern (best for 3-4 panes)
#    - If no layout is specified, use manual split/size configuration
#
# 3. PANES:
#    - Each pane is a terminal instance within a tab
#    - With predefined layouts, just list panes - layout handles arrangement
#    - With manual layout, the first pane is the base pane
#    - Subsequent panes split from the previous pane
#
# 4. MANUAL SPLITS:
#    - 'horizontal': Creates a top/bottom split (new pane appears below)
#    - 'vertical': Creates a left/right split (new pane appears to the right)
#    - If no split is specified, the pane replaces the current layout
#
# 5. SIZE:
#    - Specified as a percentage (1-100)
#    - Represents the size of the NEW pane being created
#    - For horizontal split: percentage of height
#    - For vertical split: percentage of width
#
# 6. COMMANDS:
#    - List of shell commands to run when the pane starts
#    - Commands run in order, one after another
#    - Use commands to navigate directories, start servers, open files, etc.
#
# 7. FOCUS:
#    - Only one pane per tab should have focus: true
#    - Only one tab in the entire config should have focus: true
#    - This determines what's active when the session starts
#
# TIPS:
# - Start simple with just a few tabs and panes
# - Test your configuration with: zellijinator start %s
# - You can always edit this file with: zellijinator edit %s
# - Delete a project with: zellijinator delete %s
# - List all projects with: zellijinator list
#
# For more information about Zellij: https://zellij.dev/
`, projectName, projectName, projectName, projectName, projectName, projectName, projectName, projectName, projectName)
}