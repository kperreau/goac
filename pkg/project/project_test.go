package project

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProjectsList_ValidOptions(t *testing.T) {
	opt := &Options{
		Target:         "goac",
		DryRun:         false,
		MaxConcurrency: 2,
		BinaryCheck:    true,
		Force:          false,
		DockerIgnore:   true,
		ProjectsName:   []string{"project1", "project2"},
		Debug:          []string{"debug1", "debug2"},
		PrintStdout:    true,
	}

	RootPath = "./../.."
	list, err := NewProjectsList(opt)

	assert.NoError(t, err)
	assert.NotNil(t, list)
}

func TestNewProjectsList_InvalidPath(t *testing.T) {
	opt := &Options{
		Target:         "goac",
		DryRun:         false,
		MaxConcurrency: 2,
		BinaryCheck:    true,
		Force:          false,
		DockerIgnore:   true,
		ProjectsName:   []string{"project1", "project2"},
		Debug:          []string{"debug1", "debug2"},
		PrintStdout:    true,
	}

	RootPath = "invalid-path"
	list, err := NewProjectsList(opt)

	assert.Error(t, err)
	assert.Nil(t, list)
}
