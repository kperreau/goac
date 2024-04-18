package cmd

import (
	"errors"

	"github.com/kperreau/goac/pkg/project"

	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:     "discover",
	Example: "goac discover",
	Short:   "List discovered projects",
	Long:    `Use it to discovering projects and create default config files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("bad args number")
		}

		err := project.Discover(&project.DiscoverOptions{
			Force:  force,
			Create: create,
		})
		if err != nil {
			return err
		}

		return nil
	},
}

var create bool

func init() {
	rootCmd.AddCommand(discoverCmd)

	discoverCmd.Flags().BoolVarP(&force, "force", "f", false, "Force creation file if already exist")
	discoverCmd.Flags().BoolVarP(&create, "create", "c", false, "Create project config files")
}
