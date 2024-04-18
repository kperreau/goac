package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRule_LoadsDefaultRule(t *testing.T) {
	// Initialize the project object
	p := &Project{
		Module: &Module{
			IgnoredGoFiles: []string{"file1.go", "file2.go"},
		},
	}

	// Invoke the LoadRule method
	p.LoadRule(TargetBuild)

	// Assert that the rule is loaded correctly
	assert.Equal(t, DefaultFilesToInclude[TargetBuild], p.Rule.Includes)
	assert.Equal(t, append(p.Module.IgnoredGoFiles, DefaultFilesToExclude[TargetBuild]...), p.Rule.Excludes)
}

func TestLoadRule_AppendsIgnoredGoFiles(t *testing.T) {
	// Initialize the project object
	p := &Project{
		Module: &Module{
			IgnoredGoFiles: []string{"file1.go", "file2.go"},
		},
	}

	// Invoke the LoadRule method
	p.LoadRule(TargetBuild)

	// Assert that the ignored go files are appended to the exclude list
	expectedExcludes := append(p.Module.IgnoredGoFiles, DefaultFilesToExclude[TargetBuild]...)
	assert.Equal(t, expectedExcludes, p.Rule.Excludes)
}

func TestLoadRule_HandlesEmptyIgnoredGoFiles(t *testing.T) {
	// Initialize the project object with an empty ignored go files list
	p := &Project{
		Module: &Module{
			IgnoredGoFiles: []string{},
		},
	}

	// Invoke the LoadRule method
	p.LoadRule(TargetBuild)

	// Assert that the exclude list only contains the default excludes
	assert.Equal(t, DefaultFilesToExclude[TargetBuild], p.Rule.Excludes)
}
