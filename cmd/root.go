package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goac",
	Short: "Go Affected Cache",
	Long: `GOAC is a CLI library for Go that empowers builds.
This application is a tool to check if an app is affected by recent change.
This way it improve build and deployment.`,
}

func projectsCmd(arg string) []string {
	if arg == "" {
		return []string{}
	}
	return strings.Split(arg, ",")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
