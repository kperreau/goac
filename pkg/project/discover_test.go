package project

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestDiscover_AllPotentialProjects(t *testing.T) {
	// Arrange
	opts := &DiscoverOptions{}

	// Act
	err := Discover(opts)

	// Assert
	assert.NoError(t, err)
}

func TestDiscover_PrintsNameAndPath(t *testing.T) {
	// Arrange
	opts := &DiscoverOptions{
		Create: false,
	}

	// Act
	output, err := redirectDiscoverStdout(func() error { return Discover(opts) })

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Discovered 0 potential projects\n", output.String())
}

func TestDiscover_FailedToDiscoverProjects(t *testing.T) {
	// Arrange
	opts := &DiscoverOptions{}

	// update default path to search in
	OldPahToSearch := PahToSearch
	PahToSearch = "invalid-path"
	defer func() { PahToSearch = OldPahToSearch }()

	// Act
	err := Discover(opts)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to discover projects")
}

func TestDiscover_FailedToCreateConfigFile_PermissionDenied(t *testing.T) {
	// Arrange
	opts := &DiscoverOptions{
		Create: true,
		Force:  true,
	}

	tmp, err := os.MkdirTemp("", "discover")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)
	err = os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main\nfunc main(){}"), 0o644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, configFileName), nil, 0o000)
	assert.NoError(t, err)

	// update default path to search in
	OldPahToSearch := PahToSearch
	PahToSearch = tmp
	defer func() { PahToSearch = OldPahToSearch }()

	// Act
	output, err := redirectDiscoverStdout(func() error { return Discover(opts) })

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "Failed to create project")
	assert.Contains(t, output.String(), "permission denied")
}

func TestDiscover_FailedToCreateConfigFile_AlreadyExist(t *testing.T) {
	// Arrange
	opts := &DiscoverOptions{
		Create: true,
		Force:  false,
	}

	tmp, err := os.MkdirTemp("", "discover")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)
	err = os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main\nfunc main(){}"), 0o644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, configFileName), nil, 0o644)
	assert.NoError(t, err)

	// update default path to search in
	OldPahToSearch := PahToSearch
	PahToSearch = tmp
	defer func() { PahToSearch = OldPahToSearch }()

	// Act
	output, err := redirectDiscoverStdout(func() error { return Discover(opts) })

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "[Already Exist]")
}

func TestDiscover_FailedToCreateConfigFile_Created(t *testing.T) {
	// Arrange
	opts := &DiscoverOptions{
		Create: true,
		Force:  false,
	}

	// create temp dir
	tmp, err := os.MkdirTemp("", "discover")
	assert.NoError(t, err)
	defer os.RemoveAll(tmp)
	err = os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main\nfunc main(){}"), 0o644)
	assert.NoError(t, err)

	// update default path to search in
	OldPahToSearch := PahToSearch
	PahToSearch = tmp
	defer func() { PahToSearch = OldPahToSearch }()

	// Act
	output, err := redirectDiscoverStdout(func() error { return Discover(opts) })

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, output.String(), "[Created]")
}

func TestSearchProjects_ReturnsListOfDirectories(t *testing.T) {
	// Arrange
	expected := []string{"../../main.go"}

	// update default path to search in
	OldPahToSearch := PahToSearch
	PahToSearch = "../.."
	defer func() { PahToSearch = OldPahToSearch }()

	// Act
	result, err := searchProjects()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSearchProjects_ReturnsErrorForInvalidPath(t *testing.T) {
	// Arrange
	// update default path to search in
	OldPahToSearch := PahToSearch
	PahToSearch = "invalid-path"
	defer func() { PahToSearch = OldPahToSearch }()

	// Act
	result, err := searchProjects()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPathToName_ValidPath(t *testing.T) {
	path := "/path/to/some/file"
	expected := "-path-to-some-file"

	result := pathToName(path)

	assert.Equal(t, expected, result)
}

func TestPathToName_PathIsRoot(t *testing.T) {
	path := "."
	expected := "project"

	result := pathToName(path)

	assert.Equal(t, expected, result)
}

func TestCreateConfigFile_NewFileWithGivenPathAndNameIfNotExistsAndForceIsTrue(t *testing.T) {
	// create temp dir
	tmp, err := os.MkdirTemp("", "discover")
	defer os.RemoveAll(tmp)
	assert.NoError(t, err)

	// Test setup
	path := filepath.Join(tmp, configFileName)
	name := "project-name"
	force := true

	// Execute the function
	result, err := createConfigFile(path, name, force)

	// Verify the result
	assert.NoError(t, err)
	assert.Equal(t, statusCreated, result)

	// Verify that the file was created
	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestCreateConfigFile_StatusAlreadyExistIfFileExistsAndForceIsFalse(t *testing.T) {
	// create temp dir
	tmp, err := os.MkdirTemp("", "discover")
	defer os.RemoveAll(tmp)
	assert.NoError(t, err)

	// Test setup
	path := filepath.Join(tmp, configFileName)
	name := "project-name"
	force := false

	// Create a dummy file
	file, err := os.Create(path)
	assert.NoError(t, err)
	defer func() { _ = file.Close() }()

	// Execute the function
	result, err := createConfigFile(path, name, force)

	// Verify the result
	assert.NoError(t, err)
	assert.Equal(t, statusAlreadyExist, result)
}

func TestPrintDiscoverStatus_StatusCreated(t *testing.T) {
	expected := color.GreenString("Created")
	result := printDiscoverStatus(statusCreated)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestPrintDiscoverStatus_StatusAlreadyExist(t *testing.T) {
	expected := color.HiBlackString("Already Exist")
	result := printDiscoverStatus(statusAlreadyExist)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestPrintDiscoverStatus_EmptyString(t *testing.T) {
	expected := color.RedString("Failed")
	result := printDiscoverStatus("")
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestPrintDiscoverStatus_UndefinedStatus(t *testing.T) {
	expected := color.RedString("Failed")
	result := printDiscoverStatus("undefined")
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func redirectDiscoverStdout(f func() error) (bytes.Buffer, error) {
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

	return buf, err
}
