package project

import (
	"fmt"
	"path/filepath"

	"github.com/codeskyblue/dockerignore"
	"github.com/kperreau/goac/pkg/scan"
	"github.com/kperreau/goac/pkg/utils"
)

var DefaultFilesToInclude = map[Target][]string{
	TargetBuild:      {"*.go"},
	TargetBuildImage: {},
}

var DefaultFilesToExclude = map[Target][]string{
	TargetBuild:      {".goacproject.yaml", "*_test.go"},
	TargetBuildImage: {".goacproject.yaml", "*_test.go"},
}

func (p *Project) LoadRule(target Target) {
	p.Rule = &scan.Rule{
		Includes: DefaultFilesToInclude[target],
		Excludes: append(p.Module.IgnoredGoFiles, DefaultFilesToExclude[target]...),
	}

	// add .dockerignore entries to the exclude files rules
	dockerIgnoreFiles, err := dockerignore.ReadIgnoreFile(filepath.Clean(fmt.Sprintf("%s/.dockerignore", p.CleanPath)))
	if err == nil {
		p.Rule.Excludes = utils.AppendIfNotExist(p.Rule.Excludes, dockerIgnoreFiles...)
	}
}
