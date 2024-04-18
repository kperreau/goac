package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanPath_ValidPathAndFilename(t *testing.T) {
	// Arrange
	path := "/path/to/config.yaml"
	filename := "config.yaml"

	// Act
	result := CleanPath(path, filename)

	// Assert
	assert.Equal(t, "/path/to", result)
}

func TestCleanPath_ValidDirtyPathAndFilename(t *testing.T) {
	// Arrange
	path := "./path///to//config.yaml"
	filename := "config.yaml"

	// Act
	result := CleanPath(path, filename)

	// Assert
	assert.Equal(t, "path/to", result)
}

func TestCleanPath_EmptyPathAndFilename(t *testing.T) {
	// Arrange
	path := ""
	filename := ""

	// Act
	result := CleanPath(path, filename)

	// Assert
	assert.Equal(t, ".", result)
}

func TestCleanPath_NoChange(t *testing.T) {
	// Arrange
	path := "/path/to/config.yaml"
	filename := "nop"

	// Act
	result := CleanPath(path, filename)

	// Assert
	assert.Equal(t, "/path/to/config.yaml", result)
}

func TestCleanPath_NoChange2(t *testing.T) {
	// Arrange
	path := "/path/to/nop"
	filename := "config.yaml"

	// Act
	result := CleanPath(path, filename)

	// Assert
	assert.Equal(t, "/path/to/nop", result)
}

func TestAddCurrentDirPrefix_InputPathStartsWithSlash_ReturnsInputPath(t *testing.T) {
	// Arrange
	path := "/test/path"

	// Act
	result := AddCurrentDirPrefix(path)

	// Assert
	assert.Equal(t, path, result)
}

func TestAddCurrentDirPrefix_InputPathStartsWithoutSlash_ReturnsDotSlash(t *testing.T) {
	// Arrange
	path := "test/path"

	// Act
	result := AddCurrentDirPrefix(path)

	// Assert
	assert.Equal(t, "./test/path", result)
}

func TestAddCurrentDirPrefix_InputPathIsEmptyString_ReturnsDotSlash(t *testing.T) {
	// Arrange
	path := ""

	// Act
	result := AddCurrentDirPrefix(path)

	// Assert
	assert.Equal(t, "./", result)
}

func TestFileExist_ValidFilePath(t *testing.T) {
	// Arrange
	filePath := "utils.go"

	// Act
	result := FileExist(filePath)

	// Assert
	assert.True(t, result)
}

func TestFileExist_BadFilePath(t *testing.T) {
	// Arrange
	filePath := "file-not-exist.go"

	// Act
	result := FileExist(filePath)

	// Assert
	assert.False(t, result)
}

func TestAppendIfNotExist_SingleElementToEmptySlice(t *testing.T) {
	// Create a new empty slice
	var slice []string

	// Call the AppendIfNotExist function with a single element
	slice = AppendIfNotExist(slice, "element")

	// Assert that the slice contains the appended element
	assert.Contains(t, slice, "element")
}

func TestAppendIfNotExist_ExistingElementToNonEmptySlice(t *testing.T) {
	// Create a new non-empty slice
	slice := []string{"existing"}

	// Call the AppendIfNotExist function with an existing element
	slice = AppendIfNotExist(slice, "existing")

	// Assert that the slice remains unchanged
	assert.Equal(t, []string{"existing"}, slice)
}

func TestAppendIfNotExist_AddElementsToNonEmptySlice(t *testing.T) {
	// Create a new non-empty slice
	slice := []string{"1", "2", "3"}

	// Call the AppendIfNotExist function with elements
	slice = AppendIfNotExist(slice, []string{"3", "1", "4", ""}...)

	// Assert that the slice remains changed without duplication
	assert.Equal(t, []string{"1", "2", "3", "4"}, slice)
}
