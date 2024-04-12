package cmd

import (
	"fmt"
	"os"
	"slices"
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

var (
	concurrency int
	debug       string
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	affectedCmd.Flags().StringVar(&debug, "debug", "", "Debug files loaded/hashed")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 4, "Max Concurrency")
}

var validDebugValues = []string{"name", "includes", "excludes", "dependencies", "local", "hashed"}

func debugCommand(arg string) ([]string, error) {
	if arg == "" {
		return []string{}, nil
	}
	args := strings.Split(arg, ",")
	for _, elem := range args {
		if !slices.Contains(validDebugValues, elem) {
			return []string{}, fmt.Errorf("bad debug value: %s\nvalid values are: %s\n", elem, strings.Join(validDebugValues, ","))
		}
	}

	return args, nil
}
