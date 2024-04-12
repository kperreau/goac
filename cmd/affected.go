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

		switch target {
		case project.TargetBuild.String():
			projectsList, err := project.NewProjectsList(&project.Options{
				Path:           ".",
				Target:         project.TargetBuild,
				DryRun:         dryrun,
				MaxConcurrency: concurrency,
			})
			if err != nil {
				return err
			}
			if err := projectsList.Affected(); err != nil {
				return err
			}
			return nil
		case project.TargetBuildImage.String():
			projectsList, err := project.NewProjectsList(&project.Options{
				Path:           ".",
				Target:         project.TargetBuild,
				DryRun:         dryrun,
				MaxConcurrency: concurrency,
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

var (
	build  bool
	image  bool
	target string
	dryrun bool
	//concurrency int
)

func init() {
	rootCmd.AddCommand(affectedCmd)

	affectedCmd.Flags().BoolVarP(&image, "image", "i", false, "Build Image")
	affectedCmd.Flags().BoolVarP(&build, "build", "b", false, "Build binary")
	affectedCmd.Flags().BoolVarP(&dryrun, "DryRun", "d", false, "Dry & run")
	//affectedCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 4, "Max Concurrency")
	affectedCmd.Flags().StringVarP(&target, "target", "t", "", "Target")
}
