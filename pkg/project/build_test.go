package project

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild_WithDefaultParameters(t *testing.T) {
	tmp, _ := os.MkdirTemp("", "test-build")
	p := &Project{
		Name: "test-project",
		Path: tmp,
		Target: map[Target]*TargetConfig{
			TargetBuild: {
				Exec: &Exec{
					CMD:    "echo",
					Params: []string{"go build"},
				},
			},
		},
		CMDOptions: &Options{
			Target: TargetBuild,
		},
	}

	output, err := redirectBuildStdout(p.build)

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Building test-project...")
}

func TestBuild_PrintStdout(t *testing.T) {
	p := &Project{
		Name: "test-project",
		Target: map[Target]*TargetConfig{
			TargetBuild: {
				Exec: &Exec{
					CMD:    "echo",
					Params: []string{"should print this"},
				},
			},
		},
		CMDOptions: &Options{
			Target:      TargetBuild,
			PrintStdout: true,
		},
	}

	output, err := redirectBuildStdout(p.build)

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Building test-project...")
	assert.Contains(t, output.String(), "should print this")
}

func TestBuild_WithoutPrintStdout(t *testing.T) {
	p := &Project{
		Name: "test-project",
		Target: map[Target]*TargetConfig{
			TargetBuild: {
				Exec: &Exec{
					CMD:    "echo",
					Params: []string{"should not print this"},
				},
			},
		},
		CMDOptions: &Options{
			Target: TargetBuild,
		},
	}

	output, err := redirectBuildStdout(p.build)

	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Building test-project...")
	assert.NotContains(t, output.String(), "should not print this")
}

func TestBuild_CommandFails(t *testing.T) {
	p := &Project{
		Name: "test-project",
		Target: map[Target]*TargetConfig{
			TargetBuild: {
				Exec: &Exec{
					CMD:    "cat",
					Params: []string{"not-found"},
				},
			},
		},
		CMDOptions: &Options{
			Target: TargetBuild,
		},
	}

	output, err := redirectBuildStdout(p.build)

	assert.Error(t, err)
	assert.Contains(t, output.String(), "Building test-project...")
}

func redirectBuildStdout(f func() error) (*bytes.Buffer, error) {
	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the func
	err := f()

	// Restore stdout
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
	os.Stdout = old

	// Read from the buffer and assert the output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		log.Fatal(err)
	}

	return &buf, err
}

func TestSetEnv_SetEnvironmentVariables(t *testing.T) {
	p := &Project{
		Name: "goac",
		Path: ".",
		Target: map[Target]*TargetConfig{
			TargetBuild: {
				Envs: []Env{
					{Key: "ENV_TEST", Value: "hello"},
					{Key: "BUILD_NAME", Value: "goac"},
					{Key: "PROJECT_PATH", Value: "."},
				},
			},
		},
		CMDOptions: &Options{
			Target: TargetBuild,
		},
	}
	cmd := &exec.Cmd{}
	setEnv(p, cmd)

	expectedEnv := append(
		os.Environ(),
		"ENV_TEST=hello",
		"BUILD_NAME=goac",
		"PROJECT_PATH=.",
	)
	assert.Equal(t, expectedEnv, cmd.Env)
}

func TestReplaceAllVariables_ReplaceVariables(t *testing.T) {
	p := &Project{
		Name: "goac",
		Path: ".",
		Target: map[Target]*TargetConfig{
			TargetBuild: {
				Envs: []Env{
					{Key: "PROJECT_NAME", Value: "{{project-name}}"},
					{Key: "PROJECT_PATH", Value: "{{project-path}}"},
				},
				Exec: &Exec{
					CMD:    "echo",
					Params: []string{"{{project-name}}", "{{project-path}}"},
				},
			},
		},
		CMDOptions: &Options{
			Target: TargetBuild,
		},
	}

	replaceAllVariables(p)

	expectedEnvName := p.Name
	expectedEnvPath := p.Path
	assert.Equal(t, expectedEnvName, p.Target[TargetBuild].Envs[0].Value)
	assert.Equal(t, expectedEnvPath, p.Target[TargetBuild].Envs[1].Value)

	expectedParamName := p.Name
	expectedParamPath := p.Path
	assert.Equal(t, expectedParamName, p.Target[TargetBuild].Exec.Params[0])
	assert.Equal(t, expectedParamPath, p.Target[TargetBuild].Exec.Params[1])
}
