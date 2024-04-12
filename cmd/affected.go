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
		debugArgs, err := debugCommand(debug)
		if err != nil {
			return err
		}

		switch target {
		case project.TargetBuild.String():
			projectsList, err := project.NewProjectsList(&project.Options{
				Path:           ".",
				Target:         project.TargetBuild,
				DryRun:         dryrun,
				MaxConcurrency: concurrency,
				BinaryCheck:    binaryCheck,
				Force:          force,
				DockerIgnore:   dockerignore,
				Debug:          debugArgs,
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
				Target:         project.TargetBuildImage,
				DryRun:         dryrun,
				MaxConcurrency: concurrency,
				BinaryCheck:    binaryCheck,
				Force:          force,
				DockerIgnore:   dockerignore,
				Debug:          debugArgs,
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
	target       string
	dryrun       bool
	force        bool
	binaryCheck  bool
	dockerignore bool
)

func init() {
	rootCmd.AddCommand(affectedCmd)

	affectedCmd.Flags().StringVarP(&target, "target", "t", "", "Target")
	affectedCmd.Flags().BoolVar(&dockerignore, "dockerignore", true, "Read docker ignore")
	affectedCmd.Flags().BoolVar(&binaryCheck, "binarycheck", false, "Affected if binary is missing")
	affectedCmd.Flags().BoolVar(&dryrun, "dryrun", false, "Dry & run")
	affectedCmd.Flags().BoolVarP(&force, "force", "f", false, "Force build")
}
