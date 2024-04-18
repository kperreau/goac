package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectsCmd_EmptyArgument(t *testing.T) {
	arg := ""
	expected := []string{}
	result := projectsCmd(arg)
	assert.Equal(t, expected, result)
}

func TestProjectsCmd_SplitArgument(t *testing.T) {
	arg := "project1,project2,project3"
	expected := []string{"project1", "project2", "project3"}
	result := projectsCmd(arg)
	assert.Equal(t, expected, result)
}

func TestProjectsCmd_SingleProject(t *testing.T) {
	arg := "project1"
	expected := []string{"project1"}
	result := projectsCmd(arg)
	assert.Equal(t, expected, result)
}

func TestProjectsCmd_MultipleProjects(t *testing.T) {
	arg := "project1,project2,project3"
	expected := []string{"project1", "project2", "project3"}
	result := projectsCmd(arg)
	assert.Equal(t, expected, result)
}
