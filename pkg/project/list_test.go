package project

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_PrintsNumberOfProjects(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{
			{Name: "Project1", Path: "/path/to/project1"},
			{Name: "Project2", Path: "/path/to/project2"},
		},
		Options: &Options{},
	}

	// Call the List method and return output
	buf := redirectStdout(l.List)

	// Assert that the output matches the expected value
	expectedOutput := "Found 2 projects\n" +
		"Project1 => /path/to/project1\n" +
		"Project2 => /path/to/project2\n"
	assert.Equal(t, expectedOutput, buf.String())
}

func TestList_HandlesEmptyProjectsList(t *testing.T) {
	// Initialize the List object with an empty list of projects
	l := &List{
		Projects: []*Project{},
		Options:  &Options{},
	}

	// Call the List method and return output
	buf := redirectStdout(l.List)

	// Assert that the output is empty
	assert.Equal(t, "Found 0 projects\n", buf.String())
}

func TestList_HandlesProjectsWithEmptyName(t *testing.T) {
	// Initialize the List object with a project with empty name
	l := &List{
		Projects: []*Project{
			{Name: "", Path: "."},
		},
		Options: &Options{},
	}

	// Call the List method and return output
	buf := redirectStdout(l.List)

	// Assert that the output matches the expected result
	expected := "Found 1 projects\n => .\n"
	assert.Equal(t, expected, buf.String())
}

func redirectStdout(f func()) bytes.Buffer {
	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the func
	f()

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

	return buf
}
