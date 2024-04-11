package cmd

import (
	"errors"

	"github.com/kperreau/goac/pkg/project"
	"github.com/spf13/cobra"
)

// affectedCmd represents the affected command
var affectedCmd = &cobra.Command{
	Use:     "affected",
	Short:   "List affected projects",
	Long:    `List projects affected by recent changes based on GOAC cache.`,
	Example: "goac affected -t build -d",
	RunE: func(cmd *cobra.Command, args []string) error {
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			return err
		}
		dryRun, err := cmd.Flags().GetBool("DryRun")
		if err != nil {
			return err
		}

		if target == "build" {
			projectsList, err := project.NewProjectsList(&project.Options{
				Path:           ".",
				Target:         project.TargetBuild,
				DryRun:         dryRun,
				MaxConcurrency: 4,
			})
			if err != nil {
				return err
			}
			if err := projectsList.Affected(); err != nil {
				return err
			}

			return nil
		}

		return errors.New("bad argument")
	},
}

func init() {
	rootCmd.AddCommand(affectedCmd)

	affectedCmd.Flags().BoolP("DryRun", "d", false, "DryRun")
	affectedCmd.Flags().StringP("target", "t", "", "Targets")
}
