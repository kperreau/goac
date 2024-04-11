package project

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
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

func (p *Project) LoadGOModules() error {
	cmd := exec.Command("go", "list", "-json", p.GoPath)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var rawData toolData
	if err := json.Unmarshal(output, &rawData); err != nil {
		return err
	}

	localDir, extDeps := cleanDeps(&rawData, p.GoPath)

	cmd = exec.Command("go", append([]string{"list", "-f", "{{.ImportPath}}{{if not .Standard}} {{.Module.Version}}{{end}}"}, extDeps...)...)
	output, err = cmd.Output()
	if err != nil {
		return err
	}

	externalDeps := strings.Fields(string(output))

	p.Module = &Module{
		LocalDirs:      localDir,
		ExternalDeps:   externalDeps,
		IgnoredGoFiles: rawData.IgnoredGoFiles,
	}

	return nil
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
