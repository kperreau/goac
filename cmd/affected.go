package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

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
	concurrency int
	debug       string
	projects    string
)

var validDebugValues = []string{"name", "includes", "excludes", "dependencies", "local", "hashed"}

var (
	target       string
	dryrun       bool
	force        bool
	binaryCheck  bool
	dockerignore bool
	stdout       bool
)

func debugCmd(arg string) ([]string, error) {
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

func init() {
	rootCmd.AddCommand(affectedCmd)

	affectedCmd.Flags().StringVarP(&target, "target", "t", "", "Target")
	affectedCmd.Flags().BoolVar(&stdout, "stdout", false, "Print stdout of exec command")
	affectedCmd.Flags().BoolVar(&dockerignore, "dockerignore", true, "Read docker ignore")
	affectedCmd.Flags().BoolVar(&binaryCheck, "binarycheck", false, "Affected if binary is missing")
	affectedCmd.Flags().BoolVar(&dryrun, "dryrun", false, "Dry & run")
	affectedCmd.Flags().BoolVarP(&force, "force", "f", false, "Force build")
	affectedCmd.Flags().StringVarP(&projects, "projects", "p", "", "Filter by projects name")
	affectedCmd.Flags().StringVar(&debug, "debug", "", "Display some data to debug")
	affectedCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 4, "Max Concurrency")
}
