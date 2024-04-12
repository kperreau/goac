package cmd

import (
	"os"

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

var (
	concurrency int
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 4, "Max Concurrency")
}
