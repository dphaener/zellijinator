package cmd

import (
	"fmt"
	"runtime"

	"github.com/dphaener/zellijinator/internal/styles"
	"github.com/spf13/cobra"
)

// Version information set by ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Display detailed version information about zellijinator`,
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func showVersion() {
	fmt.Println(styles.Title.Render("Zellijinator"))
	fmt.Println()
	fmt.Println(styles.InfoMsg(fmt.Sprintf("Version:   %s", styles.Bold.Render(version))))
	fmt.Println(styles.InfoMsg(fmt.Sprintf("Commit:    %s", commit)))
	fmt.Println(styles.InfoMsg(fmt.Sprintf("Built:     %s", date)))
	fmt.Println(styles.InfoMsg(fmt.Sprintf("Go:        %s", runtime.Version())))
	fmt.Println(styles.InfoMsg(fmt.Sprintf("Platform:  %s/%s", runtime.GOOS, runtime.GOARCH)))
}