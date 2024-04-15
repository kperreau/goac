package project

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestPrintAffected_1AffectedAnd1Project(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{
			{Name: "Project1", Path: "/path/to/project1", CMDOptions: &Options{DryRun: true, Force: true}},
		},
		Options: &Options{
			DryRun: true,
			Force:  true,
		},
	}

	// Call the List method and return output
	buf := redirectStdoutMessage(l.printAffected)

	// Assert the printed output
	expectedOutput := "Affected: 1/1\n"
	assert.Equal(t, expectedOutput, buf.String())
}

func TestPrintAffected_0AffectedAnd1Project(t *testing.T) {
	// Initialize the List object
	opts := &Options{
		DryRun:      true,
		BinaryCheck: false,
		Force:       false,
		Target:      TargetBuild,
	}
	l := &List{
		Projects: []*Project{
			{
				Name: "Project1", Path: "/path/to/project1",
				CMDOptions: opts,
				Metadata: &Metadata{
					DependenciesHash: "a",
					DirHash:          "b",
					Date:             "",
				},
				Cache: &Cache{Target: map[Target]*Metadata{TargetBuild: {
					DependenciesHash: "a",
					DirHash:          "b",
					Date:             "",
				}}},
			},
		},
		Options: opts,
	}

	// Call the List method and return output
	buf := redirectStdoutMessage(l.printAffected)

	// Assert the printed output
	expectedOutput := "Affected: 0/1\n"
	assert.Equal(t, expectedOutput, buf.String())
}

func redirectStdoutMessage(f func()) bytes.Buffer {
	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the func
	f()

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read from the buffer and assert the output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf
}
