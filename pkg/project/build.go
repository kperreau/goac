package project

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/printer"
)

func (p *Project) build() error {
	printer.Printf("Building %s...\n", color.HiBlueString(p.Name))

	// replace variables env and params to proper values
	replaceAllVariables(p)

	var stderr bytes.Buffer
	cmd := exec.Command(p.Target[p.CMDOptions.Target].Exec.CMD, p.Target[p.CMDOptions.Target].Exec.Params...)
	setEnv(p, cmd)
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr.String())
	}

	if p.CMDOptions.PrintStdout {
		fmt.Print(string(output))
	}

	return nil
}

func setEnv(p *Project, cmd *exec.Cmd) {
	cmd.Env = os.Environ()
	for _, env := range p.Target[p.CMDOptions.Target].Envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
}

func replaceAllVariables(p *Project) {
	// init variables
	variables := map[string]string{
		"{{project-name}}": p.Name,
		"{{project-path}}": p.Path,
	}

	for i := range p.Target[p.CMDOptions.Target].Envs {
		for search, replace := range variables {
			p.Target[p.CMDOptions.Target].Envs[i].Value = strings.ReplaceAll(p.Target[p.CMDOptions.Target].Envs[i].Value, search, replace)
		}
	}

	for i := range p.Target[p.CMDOptions.Target].Exec.Params {
		for search, replace := range variables {
			p.Target[p.CMDOptions.Target].Exec.Params[i] = strings.ReplaceAll(p.Target[p.CMDOptions.Target].Exec.Params[i], search, replace)
		}
	}
}
