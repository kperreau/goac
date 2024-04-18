package project

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

type Module struct {
	LocalDirs      []string
	ExternalDeps   []string
	IgnoredGoFiles []string
}

type toolData struct {
	ImportPath string
	Module     struct {
		Path string
		Dir  string
	}
	GoFiles        []string
	IgnoredGoFiles []string
	Imports        []string
	Deps           []string
}

func (p *Project) LoadGOModules(gomod *modfile.File) error {
	cmd := exec.Command("go", "list", "-json", p.Path)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var rawData toolData
	if err := json.Unmarshal(output, &rawData); err != nil {
		return err
	}

	localDir, extDeps := cleanDeps(&rawData, p.Path)

	p.Module = &Module{
		LocalDirs:      localDir,
		ExternalDeps:   getDependencies(gomod, extDeps),
		IgnoredGoFiles: rawData.IgnoredGoFiles,
	}

	return nil
}

func loadGOModFile(path string) (*modfile.File, error) {
	data, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("error reading go.mod: %v", err)
	}

	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("error parsing go.mod: %v", err)
	}

	return modFile, nil
}

func cleanDeps(rawData *toolData, localDir string) (localDeps []string, extDeps []string) {
	localDeps = []string{filepath.Clean(localDir)}
	deps := append(rawData.Deps, rawData.Imports...)
	for _, dep := range deps {
		if strings.Contains(dep, rawData.Module.Path) {
			path := strings.Replace(dep, rawData.Module.Path, ".", 1)
			if !slices.Contains(localDeps, path) {
				localDeps = append(localDeps, path)
			}
		} else {
			extDeps = append(extDeps, dep)
		}
	}

	return localDeps, extDeps
}

func getDependencies(gomod *modfile.File, rawDeps []string) (deps []string) {
	slices.Sort(rawDeps)
	for _, dep := range rawDeps {
		if depWithVersion := findVersion(gomod.Require, dep); depWithVersion != "" {
			deps = append(deps, depWithVersion)
		}
	}
	return deps
}

func findVersion(dependencies []*modfile.Require, val string) string {
	for _, item := range dependencies {
		if strings.Contains(val, item.Mod.Path) {
			return fmt.Sprintf("%s %s", val, item.Mod.Version)
		}
	}
	return val
}
