package project

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/printer"

	"gopkg.in/yaml.v3"
)

var PahToSearch = "."

type DiscoverOptions struct {
	Force  bool
	Create bool
}

func Discover(opts *DiscoverOptions) error {
	filesPath, err := searchProjects()
	if err != nil {
		return fmt.Errorf("failed to discover projects: %w", err)
	}

	printer.Printf("Discovered %s potential projects\n", color.YellowString("%d", len(filesPath)))
	for _, filePath := range filesPath {
		path := filepath.Clean(strings.Replace(filePath, "main.go", "", 1))
		name := pathToName(path)
		if !opts.Create {
			printer.Printf("%s %s %s\n", color.BlueString(name), color.YellowString("=>"), path)
		} else {
			statusResult, err := createConfigFile(filepath.Join(path, configFileName), name, opts.Force)
			if err != nil {
				printer.Printf("Failed to create project %s %s %s | Error: %s\n", color.BlueString(name), color.YellowString("=>"), path, color.RedString(err.Error()))
				continue
			}

			printer.Printf("%s %s %s [%s]\n", color.BlueString(name), color.YellowString("=>"), path, printDiscoverStatus(statusResult))
		}
	}

	return err
}

func searchProjects() ([]string, error) {
	// create tokens
	fset := token.NewFileSet()

	var dirs []string
	err := filepath.Walk(PahToSearch, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") {
			// Parse go file
			f, parseErr := parser.ParseFile(fset, path, nil, 0)
			if parseErr != nil {
				fmt.Println(parseErr)
				return nil
			}

			// Check if main package w/ main func exist
			if f.Name.Name == "main" {
				for _, decl := range f.Decls {
					if funcDecl, ok := decl.(*ast.FuncDecl); ok {
						if funcDecl.Name.Name == "main" {
							dirs = append(dirs, path)
							break
						}
					}
				}
			}
		}
		return nil
	})

	return dirs, err
}

func pathToName(path string) string {
	name := strings.ReplaceAll(path, "/", "-")

	if name == "" || name == "." {
		dir, err := os.Getwd()
		if err != nil {
			return "unknown"
		}
		return filepath.Base(dir)
	}

	return name
}

type status string

const (
	statusCreated      status = "created"
	statusAlreadyExist status = "already-exist"
	statusFailed       status = "failed"
)

func (s status) String() string { return string(s) }

func createConfigFile(path string, name string, force bool) (status, error) {
	_, err := os.Stat(path)
	fileExist := !os.IsNotExist(err)

	if fileExist && !force {
		return statusAlreadyExist, nil
	}

	file, err := os.Create(path)
	if err != nil {
		return statusFailed, err
	}
	defer func() { _ = file.Close() }()

	// default project config
	config := Project{
		Version: "1.0",
		Name:    name,
		Target: map[Target]*TargetConfig{
			TargetBuild: {
				Exec: &Exec{
					CMD: "go",
					Params: []string{
						"build",
						"-ldflags=-s -w",
						"-o",
						"{{project-path}}/{{project-name}}",
						"{{project-path}}",
					},
				},
			},
			TargetBuildImage: {
				Envs: []Env{
					{Key: "PROJECT_PATH", Value: "{{project-path}}"},
				},
				Exec: &Exec{
					CMD: "./_scripts/build-image.sh",
				},
			},
		},
	}

	// Encoder la structure Project en YAML et écrire dans le fichier
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(config); err != nil {
		return statusFailed, err
	}
	_ = encoder.Close()

	return statusCreated, nil
}

func printDiscoverStatus(s status) string {
	switch s {
	case statusCreated:
		return color.GreenString("Created")
	case statusAlreadyExist:
		return color.HiBlackString("Already Exist")
	}
	return color.RedString("Failed")
}
