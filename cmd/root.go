package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "zellijinator [project]",
	Short: "A CLI tool to manage Zellij sessions",
	Long: `Zellijinator is a CLI tool to manage Zellij sessions with pre-configured layouts.
	
Similar to tmuxinator for tmux, zellijinator allows you to:
- Create and manage session configurations
- Start sessions with pre-defined layouts
- List, edit, and delete session configurations`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// If a project name is provided, start it
		if len(args) == 1 {
			StartProject(args[0])
		} else {
			// No project specified, show interactive selection
			selected, err := selectProject("Select a project to start:")
			if err != nil {
				// If selection failed, show help
				cmd.Help()
				return
			}
			StartProject(selected)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zellijinator/config.yaml)")
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func initConfig() {
	// Config initialization will be implemented later
}