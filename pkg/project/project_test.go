package project

import (
	"crypto/sha1"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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

	// update default path to search in
	OldRootPath := RootPath
	RootPath = "./../.."
	defer func() { RootPath = OldRootPath }()

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

	// update default path to search in
	OldRootPath := RootPath
	RootPath = "invalid-path"
	defer func() { RootPath = OldRootPath }()

	list, err := NewProjectsList(opt)

	assert.Error(t, err)
	assert.Nil(t, list)
}

func TestFind_ValidPathAndFileName_ReturnsMatchingFilePaths(t *testing.T) {
	path := "../../"
	projectFileName := configFileName

	files, err := find(path, projectFileName)

	// Expected
	excpected := "../../.goacproject.yaml"

	assert.NoError(t, err)
	assert.Equal(t, []string{excpected}, files)
}

func TestFind_ValidPathAndFileName_NoMatchingFilePaths_ReturnsEmptyList(t *testing.T) {
	path := "."
	projectFileName := configFileName

	files, err := find(path, projectFileName)

	assert.NoError(t, err)
	assert.Empty(t, files)
}

func TestFind_EmptyPath_ReturnsError(t *testing.T) {
	path := ""
	projectFileName := configFileName

	_, err := find(path, projectFileName)

	assert.Error(t, err)
}

func TestLoadConfig_ValidConfig_ReturnsProjectObject(t *testing.T) {
	// Arrange
	file := filepath.Join("../..", configFileName)
	opts := &processProjectOptions{
		hashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
		Options: &Options{},
	}

	// Act
	project, err := loadConfig(file, opts)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "goac", project.Name)
}

func TestLoadConfig_SetCleanPathAndPathFields(t *testing.T) {
	// Arrange
	file := filepath.Join("../..", configFileName)
	opts := &processProjectOptions{
		hashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
		Options: &Options{},
	}

	// Act
	project, err := loadConfig(file, opts)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "../..", project.CleanPath)
	assert.Equal(t, "./../..", project.Path)
}

func TestLoadConfig_ErrorOpeningConfigFile_ReturnsError(t *testing.T) {
	// Arrange
	file := "nonexistent_config.yaml"
	opts := &processProjectOptions{
		hashPool: &sync.Pool{},
		Options:  &Options{},
	}

	// Act
	project, err := loadConfig(file, opts)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestLoadConfig_ErrorUnmarshallingConfigFile_ReturnsError(t *testing.T) {
	// Arrange
	file := "invalid_config.yaml"
	opts := &processProjectOptions{
		hashPool: &sync.Pool{},
		Options:  &Options{},
	}

	// Act
	project, err := loadConfig(file, opts)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestGetProjects_ValidOptionsAndProjectFilesFound(t *testing.T) {
	opt := &Options{
		Target:         TargetNone,
		DryRun:         false,
		MaxConcurrency: 5,
		BinaryCheck:    true,
		Force:          false,
		DockerIgnore:   true,
		ProjectsName:   []string{},
		Debug:          []string{},
		PrintStdout:    false,
	}

	// update default path to search in
	OldRootPath := RootPath
	RootPath = "../.."
	defer func() { RootPath = OldRootPath }()

	projects, err := getProjects(opt)

	assert.NoError(t, err)
	assert.NotNil(t, projects)
	assert.NotEmpty(t, projects)
}

func TestGetProjects_NoProjectFilesFound(t *testing.T) {
	opt := &Options{
		Target:         TargetNone,
		DryRun:         false,
		MaxConcurrency: 5,
		BinaryCheck:    true,
		Force:          false,
		DockerIgnore:   true,
		ProjectsName:   []string{},
		Debug:          []string{},
		PrintStdout:    false,
	}

	projects, err := getProjects(opt)

	assert.Error(t, err)
	assert.Nil(t, projects)
}

func TestProcessProject_ProjectFound(t *testing.T) {
	// Arrange
	sem := make(chan bool, 2)
	projectsCh := make(chan *Project)
	errorsCh := make(chan error)
	projectFile := filepath.Join("../..", configFileName)
	wg := sync.WaitGroup{}
	mfile, err := loadGOModFile("./../..")
	assert.NoError(t, err)
	opts := &processProjectOptions{
		hashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
		Options: &Options{
			Target: TargetNone,
		},
		sem:       sem,
		projectCh: projectsCh,
		errorsCh:  errorsCh,
		wg:        &wg,
		gomod:     mfile,
	}

	// Act
	sem <- true // acquire
	wg.Add(1)
	go processProject(opts, projectFile)
	wg.Wait()

	var errorsProjects error
	var project *Project
	select {
	case project = <-projectsCh:
	case errorsProjects = <-errorsCh:
	}

	// Assert
	assert.NoError(t, errorsProjects)
	assert.NotEmpty(t, project)
	assert.Equal(t, "goac", project.Name)
}

func TestProcessProject_SkipProjectIfProjectsOptionNotMatch(t *testing.T) {
	// Arrange
	sem := make(chan bool, 3)
	projectsCh := make(chan *Project)
	errorsCh := make(chan error)
	projectFile := filepath.Join("../..", configFileName)
	wg := sync.WaitGroup{}
	opts := &processProjectOptions{
		hashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
		Options: &Options{
			ProjectsName: []string{"project1", "project2"},
		},
		sem:       sem,
		projectCh: projectsCh,
		errorsCh:  errorsCh,
		wg:        &wg,
	}

	// Act
	sem <- true // acquire
	wg.Add(1)
	go processProject(opts, projectFile)
	wg.Wait()

	var errorsProjects error
	var project *Project
	select {
	case project = <-projectsCh:
	case errorsProjects = <-errorsCh:
	}

	// Assert
	assert.Empty(t, project)
	assert.NoError(t, errorsProjects)
}

func TestProcessProject_NoProjectFilesFound(t *testing.T) {
	// Arrange
	sem := make(chan bool, 3)
	projectsCh := make(chan *Project)
	errorsCh := make(chan error)
	projectFile := filepath.Join("invalid-path", configFileName)
	wg := sync.WaitGroup{}
	opts := &processProjectOptions{
		hashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
		Options: &Options{
			ProjectsName: []string{"goac"},
		},
		sem:       sem,
		projectCh: projectsCh,
		errorsCh:  errorsCh,
		wg:        &wg,
	}

	// Act
	sem <- true // acquire
	wg.Add(1)
	go processProject(opts, projectFile)
	wg.Wait()

	var err error
	var project *Project
	select {
	case project = <-projectsCh:
	case err = <-errorsCh:
	}

	// Assert
	assert.Empty(t, project)
	assert.Error(t, err)
}
