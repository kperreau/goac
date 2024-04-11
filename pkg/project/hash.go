package project

import (
	"encoding/hex"
	"hash"
	"strings"

	"github.com/kperreau/goac/pkg/hasher"
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

	h := p.hashPool.Get().(hash.Hash)
	defer p.hashPool.Put(h)
	h.Reset()

	if _, err := h.Write([]byte(joinedDeps)); err != nil {
		return "", err
	}

	hashBytes := h.Sum(nil)
	hashStr := hex.EncodeToString(hashBytes)

	return hashStr, nil
}

func processDirectoryHash(p *Project) (string, error) {
	files, err := scan.Dirs(p.Module.LocalDeps, p.Rule)
	if err != nil {
		return "", err
	}

	// TODO: make optional cli --debug to print local files match
	// fmt.Println("NAME", p.Name)
	// fmt.Println("")
	// fmt.Println(strings.Join(files, "\n"))
	// fmt.Println("")

	hashDir, err := hasher.Files(files, p.hashPool)
	if err != nil {
		return "", err
	}

	return hashDir, nil
}
