package project

import (
	"os/exec"
	"path"

	"github.com/kperreau/goac/pkg/printer"
)

func (p *Project) buildProject() (string, error) {
	printer.Printf("Building %s...\n", MessageName(p.Name))
	cmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", path.Join(p.GoPath, p.Name), p.GoPath)
	output, err := cmd.Output()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
