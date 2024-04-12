package cmd

import (
	"errors"

	"github.com/kperreau/goac/pkg/project"

	"github.com/spf13/cobra"
)

// listCmd represents the project command
var listCmd = &cobra.Command{
	Use:     "list",
	Example: "goac list",
	Short:   "List projects",
	Long:    `Use it to list all your projects configured with GOAC.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("bad args number")
		}

		debugArgs, err := debugCommand(debug)
		if err != nil {
			return err
		}

		listProject, err := project.NewProjectsList(&project.Options{
			Path:           ".",
			Target:         project.TargetNone,
			MaxConcurrency: concurrency,
			BinaryCheck:    binaryCheck,
			DockerIgnore:   dockerignore,
			Debug:          debugArgs,
		})
		if err != nil {
			return err
		}

		listProject.List()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
