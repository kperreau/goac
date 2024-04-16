package project

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/printer"
)

func (p *Project) build() error {
	printer.Printf("Building %s...\n", color.HiBlueString(p.Name))
	params := append(p.Target[p.CMDOptions.Target].Exec.Params, path.Join(p.GoPath, p.Name), p.GoPath)

	var stderr bytes.Buffer
	cmd := exec.Command(p.Target[p.CMDOptions.Target].Exec.CMD, params...)
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
	cmd.Env = append(os.Environ(), fmt.Sprintf("BUILD_NAME=%s", p.Name))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PROJECT_PATH=%s", p.GoPath))
}
