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
	Example: "goac affected -t build",
	RunE: func(cmd *cobra.Command, args []string) error {
		debugArgs, err := debugCmd(debug)
		if err != nil {
			return err
		}

		t := project.StringToTarget(target)
		if project.StringToTarget(target) != project.TargetNone {
			projectsList, err := project.NewProjectsList(&project.Options{
				Path:           ".",
				Target:         t,
				DryRun:         dryrun,
				MaxConcurrency: concurrency,
				BinaryCheck:    binaryCheck,
				Force:          force,
				DockerIgnore:   dockerignore,
				Debug:          debugArgs,
				ProjectsName:   projectsCmd(projects),
				PrintStdout:    stdout,
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
	stdout       bool
)

func init() {
	rootCmd.AddCommand(affectedCmd)

	affectedCmd.Flags().StringVarP(&target, "target", "t", "", "Target")
	affectedCmd.Flags().BoolVar(&stdout, "stdout", false, "Print stdout of exec command")
	affectedCmd.Flags().BoolVar(&dockerignore, "dockerignore", true, "Read docker ignore")
	affectedCmd.Flags().BoolVar(&binaryCheck, "binarycheck", false, "Affected if binary is missing")
	affectedCmd.Flags().BoolVar(&dryrun, "dryrun", false, "Dry & run")
	affectedCmd.Flags().BoolVarP(&force, "force", "f", false, "Force build")
}
