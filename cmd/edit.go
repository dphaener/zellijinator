package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dphaener/zellijinator/config"
	"github.com/dphaener/zellijinator/internal/styles"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [project]",
	Short: "Edit a zellijinator project",
	Long:  `Open the specified project configuration in your default editor`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string
		
		if len(args) == 0 {
			// No project specified, show interactive selection
			selected, err := selectProject("Select a project to edit:")
			if err != nil {
				fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("%v", err)))
				os.Exit(1)
			}
			projectName = selected
		} else {
			projectName = args[0]
		}
		
		editProject(projectName)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func editProject(name string) {
	// Get project path
	projectPath := config.ProjectPath(name)
	
	// Check if project exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Project '%s' not found", name)))
		fmt.Fprintln(os.Stderr, styles.InfoMsg(fmt.Sprintf("Create it with: %s", styles.Command.Render(fmt.Sprintf("zellijinator new %s", name)))))
		os.Exit(1)
	}
	
	// Open in editor (reuse the openInEditor function from new.go)
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
		fmt.Println(styles.InfoMsg("You can manually edit: " + styles.Path.Render(projectPath)))
		return
	}

	fmt.Println(styles.InfoMsg(fmt.Sprintf("Opening %s in %s...", styles.Bold.Render(name), styles.Command.Render(editor))))
	editorCmd := exec.Command(editor, projectPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	
	if err := editorCmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, styles.ErrorMsg(fmt.Sprintf("Error opening editor: %v", err)))
		fmt.Println(styles.InfoMsg("You can manually edit: " + styles.Path.Render(projectPath)))
	}
}