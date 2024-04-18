package project

import (
	"encoding/hex"
	"hash"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/hasher"
	"github.com/kperreau/goac/pkg/printer"
	"github.com/kperreau/goac/pkg/scan"
)

type Metadata struct {
	DependenciesHash string
	DirHash          string
	Date             string
}

func (p *Project) LoadHashs() error {
	depsHash, err := processDependenciesHash(p)
	if err != nil {
		return err
	}

	dirHash, err := processDirectoryHash(p)
	if err != nil {
		return err
	}

	p.Metadata = &Metadata{
		DependenciesHash: depsHash,
		DirHash:          dirHash,
	}

	return nil
}

func processDependenciesHash(p *Project) (string, error) {
	joinedDeps := strings.Join(p.Module.ExternalDeps, ",")

	h := p.HashPool.Get().(hash.Hash)
	defer p.HashPool.Put(h)
	h.Reset()

	if _, err := h.Write([]byte(joinedDeps)); err != nil {
		return "", err
	}

	hashBytes := h.Sum(nil)
	hashStr := hex.EncodeToString(hashBytes)

	return hashStr, nil
}

func processDirectoryHash(p *Project) (string, error) {
	files, err := scan.Dirs(p.Module.LocalDirs, p.Rule)
	if err != nil {
		return "", err
	}

	if len(p.CMDOptions.Debug) > 0 {
		debug(p, files)
	}

	hashDir, err := hasher.Files(files, p.HashPool)
	if err != nil {
		return "", err
	}

	return hashDir, nil
}

func debug(p *Project, files []string) {
	if slices.Contains(p.CMDOptions.Debug, "name") {
		printer.Warnf("Name: %s\n", printer.BoldGreen(p.Name))
	}

	if slices.Contains(p.CMDOptions.Debug, "includes") {
		printer.Printf("%s\n%s\n", color.YellowString("Includes"), strings.Join(p.Rule.Includes, "\n"))
	}

	if slices.Contains(p.CMDOptions.Debug, "excludes") {
		printer.Printf("%s\n%s\n", color.YellowString("Excludes"), strings.Join(p.Rule.Excludes, "\n"))
	}

	if slices.Contains(p.CMDOptions.Debug, "hashed") {
		printer.Printf("%s\n%s\n", color.YellowString("Hashed files"), strings.Join(files, "\n"))
	}

	if slices.Contains(p.CMDOptions.Debug, "dependencies") {
		printer.Printf("%s\n%s\n", color.YellowString("Dependencies"), strings.Join(p.Module.ExternalDeps, "\n"))
	}

	if slices.Contains(p.CMDOptions.Debug, "local") {
		printer.Printf("%s\n%s\n", color.YellowString("Local Imports"), strings.Join(p.Module.LocalDirs, "\n"))
	}

	if len(p.CMDOptions.Debug) > 0 {
		printer.Printf("\n")
	}
}
