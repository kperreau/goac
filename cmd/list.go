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

		listProject, err := project.NewProjectsList(&project.Options{
			Target:         project.TargetNone,
			MaxConcurrency: concurrency,
			ProjectsName:   projectsCmd(projects),
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

	listCmd.Flags().StringVarP(&projects, "projects", "p", "", "Filter by projects name")
	listCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 4, "Max Concurrency")
}
