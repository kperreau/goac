package scan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirs_ReturnsListOfFiles(t *testing.T) {
	// Create a mock rule
	rule := &Rule{
		Excludes: []string{},
		Includes: []string{},
	}

	// Set Directories
	dirs := []string{"."}

	// Call the Dirs function
	files, err := Dirs(dirs, rule)

	// Assert that the function returns the expected result
	assert.Nil(t, err)
	assert.Equal(t, []string{"scan.go", "scan_test.go"}, files)
}

func TestDirs_ReturnsErrorWhenDirectoryDoesNotExist(t *testing.T) {
	// Create a mock rule
	rule := &Rule{
		Excludes: []string{},
		Includes: []string{},
	}

	// Directory that does not exist
	dirs := []string{"/path/to/nonexistent/directory"}

	// Call the Dirs function
	files, err := Dirs(dirs, rule)

	// Assert that the function returns the expected error
	assert.NotNil(t, err)
	assert.Equal(t, "lstat /path/to/nonexistent/directory: no such file or directory", err.Error())

	// Assert that the function does not return any files
	assert.Empty(t, files)
}

func TestSubDir_ValidDirectoryAndRule_ReturnsListOfFiles(t *testing.T) {
	// Create a mock rule
	rule := &Rule{
		Excludes: []string{".*"},
		Includes: []string{"*.go"},
	}

	// Call the subDir function
	files, err := subDir(".", rule)

	// Assert that the function returns the expected result
	assert.NoError(t, err)
	assert.Equal(t, []string{"scan.go", "scan_test.go"}, files)
}

func TestSubDir_ValidDirectoryAndRuleExcludeTestFiles_ReturnsListOfFiles(t *testing.T) {
	// Create a mock rule
	rule := &Rule{
		Excludes: []string{".*", "*_test.go"},
		Includes: []string{"*.go"},
	}

	// Call the subDir function
	files, err := subDir(".", rule)

	// Assert that the function returns the expected result
	assert.NoError(t, err)
	assert.Equal(t, []string{"scan.go"}, files)
}

func TestSubDir_ValidDirectoryAndRuleOnlyGO_ReturnsListOfFiles(t *testing.T) {
	// Create a mock rule
	rule := &Rule{
		Excludes: []string{".*", "*_test.go"},
		Includes: []string{"*.go"},
	}

	// Call the subDir function
	files, err := subDir("../../", rule)

	// Assert that the function returns the expected result
	assert.NoError(t, err)
	assert.Contains(t, files, "../../main.go")
	assert.NotContains(t, files, "../../.goacproject.yaml")
	assert.NotContains(t, files, "../../go.mod")
}

func TestSubDir_InvalidDirectoryPath_ReturnsError(t *testing.T) {
	// Create a mock rule
	rule := &Rule{
		Excludes: []string{".*"},
		Includes: []string{"*.go"},
	}

	// Call the subDir function with an invalid directory path
	files, err := subDir("/invalid/path", rule)

	// Assert that the function returns an error
	assert.Error(t, err)
	assert.Nil(t, files)
}

func TestFileMatch_MatchesPattern_ReturnsTrue(t *testing.T) {
	// Arrange
	filename := "example.txt"
	patterns := []string{"*.txt"}

	// Act
	result := fileMatch(filename, patterns)

	// Assert
	assert.True(t, result)
}

func TestFileMatch_MatchesMultiplePatterns_ReturnsTrue(t *testing.T) {
	// Arrange
	filename := "scan.go"
	patterns := []string{"go", ".", "*.go"}

	// Act
	result := fileMatch(filename, patterns)

	// Assert
	assert.True(t, result)
}

func TestFileMatch_MatchesMultiplePatterns_ReturnsFalse(t *testing.T) {
	// Arrange
	filename := "scan.go.yo"
	patterns := []string{"go", ".", "*.go", "scan"}

	// Act
	result := fileMatch(filename, patterns)

	// Assert
	assert.False(t, result)
}

func TestFileMatch_EmptyFilename_ReturnsFalse(t *testing.T) {
	// Arrange
	filename := ""
	patterns := []string{"*.txt"}

	// Act
	result := fileMatch(filename, patterns)

	// Assert
	assert.False(t, result)
}
