package zellij

import (
	"fmt"
	"os"
	"strings"

	"github.com/dphaener/zellijinator/config"
)

// GenerateLayout creates a Zellij layout in KDL format from a project config
func GenerateLayout(project *config.Project) string {
	var layout strings.Builder

	// Set the session name
	sessionName := project.SessionName
	if sessionName == "" {
		sessionName = project.Name
	}
	layout.WriteString(fmt.Sprintf("session_name \"%s\"\n\n", sessionName))

	// If a default layout is specified, we need to structure it differently
	// Zellij expects the layout to extend from a base layout
	if project.DefaultLayout != "" {
		// For layouts like "compact", we extend from them
		layout.WriteString(fmt.Sprintf("layout {\n"))
		layout.WriteString(fmt.Sprintf("    // Extending from %s layout\n", project.DefaultLayout))
		
		// Import the compact layout's tab template
		if project.DefaultLayout == "compact" {
			layout.WriteString("    default_tab_template {\n")
			layout.WriteString("        children\n")
			layout.WriteString("        pane size=1 borderless=true {\n")
			layout.WriteString("            plugin location=\"zellij:compact-bar\"\n")
			layout.WriteString("        }\n")
			layout.WriteString("    }\n\n")
		}
	} else {
		layout.WriteString("layout {\n")
		// Standard layout with full plugins
		layout.WriteString("    default_tab_template {\n")
		layout.WriteString("        pane size=1 borderless=true {\n")
		layout.WriteString("            plugin location=\"zellij:tab-bar\"\n")
		layout.WriteString("        }\n")
		layout.WriteString("        children\n")
		layout.WriteString("        pane size=2 borderless=true {\n")
		layout.WriteString("            plugin location=\"zellij:status-bar\"\n")
		layout.WriteString("        }\n")
		layout.WriteString("    }\n\n")
	}

	// Find the focused tab
	focusedTabIndex := 0
	for i, tab := range project.Tabs {
		if tab.Focus {
			focusedTabIndex = i
			break
		}
	}

	// Generate tabs
	for tabIndex, tab := range project.Tabs {
		isFocusedTab := tabIndex == focusedTabIndex

		layout.WriteString(fmt.Sprintf("    tab name=\"%s\"", tab.Name))
		
		// Add cwd on the same line as tab declaration
		layout.WriteString(fmt.Sprintf(" cwd=\"%s\"", project.Root))
		
		if isFocusedTab {
			layout.WriteString(" focus=true")
		}
		layout.WriteString(" {\n")

		// Generate panes for this tab
		if tab.Layout != "" {
			// Use predefined layout
			generatePanesWithLayout(&layout, tab.Layout, tab.Panes, project.Root, project.Env, "        ")
		} else {
			// Use manual layout from pane definitions
			generatePanes(&layout, tab.Panes, project.Root, project.Env, "        ")
		}

		layout.WriteString("    }\n")
	}

	layout.WriteString("}\n")

	return layout.String()
}

// Helper function to write a split container with a pane inside
func writeSplitPane(layout *strings.Builder, pane *config.Pane, splitDir string, size string, rootDir string, envVars map[string]string, indent string) {
	layout.WriteString(fmt.Sprintf("%spane split_direction=\"%s\"", indent, splitDir))
	if size != "" {
		layout.WriteString(fmt.Sprintf(" size=\"%s\"", size))
	}
	layout.WriteString(" {\n")
	
	layout.WriteString(fmt.Sprintf("%s    pane", indent))
	if pane.Focus {
		layout.WriteString(" focus=true")
	}
	layout.WriteString(" {\n")
	writePaneCommand(layout, pane, rootDir, envVars, indent+"        ")
	layout.WriteString(fmt.Sprintf("%s    }\n", indent))
	
	layout.WriteString(fmt.Sprintf("%s}\n", indent))
}

// generatePanesWithLayout generates panes using a predefined layout pattern
func generatePanesWithLayout(layout *strings.Builder, layoutType string, panes []config.Pane, rootDir string, envVars map[string]string, indent string) {
	numPanes := len(panes)
	if numPanes == 0 {
		return
	}

	switch layoutType {
	case "even-horizontal":
		// All panes split horizontally with equal size
		size := 100 / numPanes
		for i, pane := range panes {
			if i == 0 {
				writePaneWithCommand(layout, &pane, rootDir, envVars, indent)
			} else {
				writeSplitPane(layout, &pane, "horizontal", fmt.Sprintf("%d%%", size), rootDir, envVars, indent)
			}
		}

	case "even-vertical":
		// All panes split vertically with equal size
		size := 100 / numPanes
		for i, pane := range panes {
			if i == 0 {
				writePaneWithCommand(layout, &pane, rootDir, envVars, indent)
			} else {
				writeSplitPane(layout, &pane, "vertical", fmt.Sprintf("%d%%", size), rootDir, envVars, indent)
			}
		}

	case "main-vertical":
		// First pane takes 70%, others split the remaining 30% horizontally
		if numPanes == 1 {
			writePaneWithCommand(layout, &panes[0], rootDir, envVars, indent)
		} else {
			// First pane
			writePaneWithCommand(layout, &panes[0], rootDir, envVars, indent)
			
			// Create vertical split container for remaining panes
			layout.WriteString(fmt.Sprintf("%spane split_direction=\"vertical\" size=\"30%%\" {\n", indent))
			
			// Second pane
			pane := panes[1]
			layout.WriteString(fmt.Sprintf("%s    pane", indent))
			if pane.Focus {
				layout.WriteString(" focus=true")
			}
			layout.WriteString(" {\n")
			writePaneCommand(layout, &pane, rootDir, envVars, indent+"        ")
			layout.WriteString(fmt.Sprintf("%s    }\n", indent))
			
			// Remaining panes split horizontally within the 30%
			if numPanes > 2 {
				remainingSize := 100 / (numPanes - 1)
				for i := 2; i < numPanes; i++ {
					writeSplitPane(layout, &panes[i], "horizontal", fmt.Sprintf("%d%%", remainingSize), rootDir, envVars, indent+"    ")
				}
			}
			layout.WriteString(fmt.Sprintf("%s}\n", indent))
		}

	case "main-horizontal":
		// First pane takes 70%, others split the remaining 30% vertically
		if numPanes == 1 {
			writePaneWithCommand(layout, &panes[0], rootDir, envVars, indent)
		} else {
			// First pane
			writePaneWithCommand(layout, &panes[0], rootDir, envVars, indent)
			
			// Create horizontal split container for remaining panes
			layout.WriteString(fmt.Sprintf("%spane split_direction=\"horizontal\" size=\"30%%\" {\n", indent))
			
			// Second pane
			pane := panes[1]
			layout.WriteString(fmt.Sprintf("%s    pane", indent))
			if pane.Focus {
				layout.WriteString(" focus=true")
			}
			layout.WriteString(" {\n")
			writePaneCommand(layout, &pane, rootDir, envVars, indent+"        ")
			layout.WriteString(fmt.Sprintf("%s    }\n", indent))
			
			// Remaining panes split vertically within the 30%
			if numPanes > 2 {
				remainingSize := 100 / (numPanes - 1)
				for i := 2; i < numPanes; i++ {
					writeSplitPane(layout, &panes[i], "vertical", fmt.Sprintf("%d%%", remainingSize), rootDir, envVars, indent+"    ")
				}
			}
			layout.WriteString(fmt.Sprintf("%s}\n", indent))
		}

	case "tiled":
		// For tiled layout, we try to create a grid
		// This is more complex - for now, let's do a simple version
		if numPanes <= 2 {
			// Just split vertically for 2 panes
			generatePanesWithLayout(layout, "even-vertical", panes, rootDir, envVars, indent)
		} else if numPanes == 3 {
			// One on left, two on right
			writePaneWithCommand(layout, &panes[0], rootDir, envVars, indent)
			
			// Right side container
			layout.WriteString(fmt.Sprintf("%spane split_direction=\"vertical\" size=\"50%%\" {\n", indent))
			
			// Top right pane
			layout.WriteString(fmt.Sprintf("%s    pane", indent))
			if panes[1].Focus {
				layout.WriteString(" focus=true")
			}
			layout.WriteString(" {\n")
			writePaneCommand(layout, &panes[1], rootDir, envVars, indent+"        ")
			layout.WriteString(fmt.Sprintf("%s    }\n", indent))
			
			// Bottom right pane
			writeSplitPane(layout, &panes[2], "horizontal", "50%", rootDir, envVars, indent+"    ")
			
			layout.WriteString(fmt.Sprintf("%s}\n", indent))
		} else if numPanes == 4 {
			// 2x2 grid
			// Top left
			writePaneWithCommand(layout, &panes[0], rootDir, envVars, indent)
			
			// Top right
			writeSplitPane(layout, &panes[1], "vertical", "50%", rootDir, envVars, indent)
			
			// Bottom row container
			layout.WriteString(fmt.Sprintf("%spane split_direction=\"horizontal\" size=\"50%%\" {\n", indent))
			
			// Bottom left
			layout.WriteString(fmt.Sprintf("%s    pane", indent))
			if panes[2].Focus {
				layout.WriteString(" focus=true")
			}
			layout.WriteString(" {\n")
			writePaneCommand(layout, &panes[2], rootDir, envVars, indent+"        ")
			layout.WriteString(fmt.Sprintf("%s    }\n", indent))
			
			// Bottom right
			writeSplitPane(layout, &panes[3], "vertical", "50%", rootDir, envVars, indent+"    ")
			
			layout.WriteString(fmt.Sprintf("%s}\n", indent))
		} else {
			// For more than 4, fall back to even-horizontal
			generatePanesWithLayout(layout, "even-horizontal", panes, rootDir, envVars, indent)
		}

	default:
		// Unknown layout, fall back to manual
		generatePanes(layout, panes, rootDir, envVars, indent)
	}
}

func generatePanes(layout *strings.Builder, panes []config.Pane, rootDir string, envVars map[string]string, indent string) {
	if len(panes) == 0 {
		return
	}

	// For multiple panes, we need to structure them properly
	if len(panes) > 1 {
		// First pane - no split needed
		writePaneWithCommand(layout, &panes[0], rootDir, envVars, indent)
		
		// For subsequent panes, we need to create container panes with splits
		for i := 1; i < len(panes); i++ {
			pane := panes[i]
			
			// Create a container pane with split direction
			if pane.Split == "horizontal" {
				layout.WriteString(fmt.Sprintf("%spane split_direction=\"horizontal\"", indent))
			} else if pane.Split == "vertical" {
				layout.WriteString(fmt.Sprintf("%spane split_direction=\"vertical\"", indent))
			} else {
				// Default to horizontal if not specified
				layout.WriteString(fmt.Sprintf("%spane split_direction=\"horizontal\"", indent))
			}
			
			if pane.Size != "" {
				layout.WriteString(fmt.Sprintf(" size=\"%s%%\"", pane.Size))
			}
			
			layout.WriteString(" {\n")
			
			// Inside the container, create the actual pane with the command
			layout.WriteString(fmt.Sprintf("%s    pane", indent))
			if pane.Focus {
				layout.WriteString(" focus=true")
			}
			layout.WriteString(" {\n")
			writePaneCommand(layout, &pane, rootDir, envVars, indent+"        ")
			layout.WriteString(fmt.Sprintf("%s    }\n", indent))
			
			layout.WriteString(fmt.Sprintf("%s}\n", indent))
		}
	} else {
		// Single pane
		writePaneWithCommand(layout, &panes[0], rootDir, envVars, indent)
	}
}

func writePaneWithCommand(layout *strings.Builder, pane *config.Pane, rootDir string, envVars map[string]string, indent string) {
	layout.WriteString(fmt.Sprintf("%spane", indent))
	
	if pane.Focus {
		layout.WriteString(" focus=true")
	}
	
	layout.WriteString(" {\n")
	writePaneCommand(layout, pane, rootDir, envVars, indent+"    ")
	layout.WriteString(fmt.Sprintf("%s}\n", indent))
}

func writePaneCommand(layout *strings.Builder, pane *config.Pane, rootDir string, envVars map[string]string, indent string) {
	// Get the user's default shell from SHELL environment variable
	userShell := os.Getenv("SHELL")
	if userShell == "" {
		userShell = "/bin/bash" // fallback to bash if SHELL is not set
	}
	
	if len(pane.Commands) > 0 {
		// Prepare environment variables
		envStr := ""
		if len(envVars) > 0 {
			var envPairs []string
			for k, v := range envVars {
				envPairs = append(envPairs, fmt.Sprintf("export %s='%s'", k, v))
			}
			envStr = strings.Join(envPairs, "; ") + "; "
		}
		
		// Join all commands
		commandStr := strings.Join(pane.Commands, " && ")
		
		// Create the full command with environment and directory change
		// Use the user's shell and exec to it after commands complete
		fullCommand := fmt.Sprintf("%scd '%s' && %s; exec %s", envStr, rootDir, commandStr, userShell)
		
		// Escape quotes in command
		fullCommand = strings.ReplaceAll(fullCommand, `"`, `\"`)
		
		// Use sh to run the command (more portable than bash)
		layout.WriteString(fmt.Sprintf("%scommand \"sh\"\n", indent))
		layout.WriteString(fmt.Sprintf("%sargs \"-c\" \"%s\"\n", indent, fullCommand))
	} else {
		// If no commands, just start the user's shell in the right directory
		fullCommand := fmt.Sprintf("cd '%s'; exec %s", rootDir, userShell)
		layout.WriteString(fmt.Sprintf("%scommand \"sh\"\n", indent))
		layout.WriteString(fmt.Sprintf("%sargs \"-c\" \"%s\"\n", indent, fullCommand))
	}
}