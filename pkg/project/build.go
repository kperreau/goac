package project

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/printer"
	"os"
	"os/exec"
	"path"
)

func (p *Project) build() (string, error) {
	printer.Printf("Building %s...\n", color.HiBlueString(p.Name))

	params := append(p.Target[p.CMDOptions.Target].Exec.Params, path.Join(p.GoPath, p.Name), p.GoPath)

	var stderr bytes.Buffer
	cmd := exec.Command(p.Target[p.CMDOptions.Target].Exec.CMD, params...)
	setEnv(p, cmd)
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return string(output), err
	}

	if p.CMDOptions.PrintStdout {
		fmt.Print(string(output))
	}

	return string(output), nil
}

func setEnv(p *Project, cmd *exec.Cmd) {
	cmd.Env = append(os.Environ(), fmt.Sprintf("BUILD_NAME=%s", p.Name))
	if p.CMDOptions.Target == TargetBuildImage {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PROJECT_PATH=%s", p.GoPath))
	}
}
