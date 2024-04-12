package project

import (
	"github.com/kperreau/goac/pkg/printer"
	"os/exec"
	"path"
)

func (p *Project) build() (string, error) {
	printer.Printf("Building %s...\n", MessageName(p.Name))

	params := append(p.Target[p.CMDOptions.Target].Exec.Params, path.Join(p.GoPath, p.Name), p.GoPath)
	cmd := exec.Command(p.Target[p.CMDOptions.Target].Exec.CMD, params...)
	output, err := cmd.Output()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
