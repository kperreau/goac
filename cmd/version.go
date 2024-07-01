package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// versionCmd represents the project command
var versionCmd = &cobra.Command{
	Use:     "version",
	Example: "goac version",
	Short:   "Get goac version",
	Long:    `Use it to know your goac version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("bad args number")
		}

		return run()
	},
}

type Semver struct {
	Alpha   int    `yaml:"alpha"`
	Beta    int    `yaml:"beta"`
	RC      int    `yaml:"rc"`
	Release string `yaml:"release"`
}

func run() error {
	// read file .semver.yaml
	data, err := os.ReadFile(".semver.yaml")
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// unmarshall YAML
	var semver Semver
	err = yaml.Unmarshal(data, &semver)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// Print 'release'
	fmt.Printf("goac version %s\n", semver.Release)

	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
