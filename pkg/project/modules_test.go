package project

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

func TestLoadGOModules_LoadsModuleData(t *testing.T) {
	// Initialize the class object
	p := &Project{
		Name: "test-project",
		Path: "./../..",
	}

	mfile, _ := loadGOModFile("../..")

	// Call the method under test
	err := p.LoadGOModules(mfile)

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the module data is correctly parsed
	assert.Contains(t, p.Module.LocalDirs, "../..")
	assert.Contains(t, p.Module.LocalDirs, "./cmd")
	assert.Contains(t, p.Module.LocalDirs, "./pkg/project")
	assert.Contains(t, p.Module.LocalDirs, "./pkg/hasher")
	assert.NotEmpty(t, p.Module.ExternalDeps)
}

func TestLoadGOModFile_Valid(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a valid go.mod file
	goModContent := "module example.com\n\ngo 1.22"
	goModPath := filepath.Join(tempDir, "go.mod")
	err := os.WriteFile(goModPath, []byte(goModContent), 0o644)
	assert.NoError(t, err)

	// Call the loadGOModFile function
	modFile, err := loadGOModFile(tempDir)
	assert.NoError(t, err)

	// Assert that the modFile is not empty
	assert.NotEmpty(t, modFile)
}

func TestLoadGOModFile_GOAC(t *testing.T) {
	// Call the loadGOModFile function
	modFile, err := loadGOModFile("../..")
	assert.NoError(t, err)

	// Assert that the modFile is not empty
	assert.NotEmpty(t, modFile)
}

func TestLoadGOModFile_LocalDependencies(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a go.mod file with local dependencies
	goModContent := "module example.com\n\ngo 1.22\n\nrequire (\n\texample.com/localdep v1.0.0\n)"
	goModPath := filepath.Join(tempDir, "go.mod")
	err := os.WriteFile(goModPath, []byte(goModContent), 0o644)
	assert.NoError(t, err)

	// Call the loadGOModFile function
	modFile, err := loadGOModFile(tempDir)
	assert.NoError(t, err)

	// Assert that the modFile is not nil
	assert.NotEmpty(t, modFile)

	// Assert that the Require slice is not empty
	assert.NotEmpty(t, modFile.Require)
	assert.Len(t, modFile.Require, 1)
}

func TestLoadGOModFile_InvalidPath(t *testing.T) {
	// Call the loadGOModFile function
	modFile, err := loadGOModFile("invalid")

	// assert error and modfile empty
	assert.Error(t, err)
	assert.Empty(t, modFile)
}

func TestCleanDeps_IncludesAllDeps(t *testing.T) {
	rawData := &toolData{
		Module: struct {
			Path string
			Dir  string
		}{
			Path: "github.com/kperreau/goac",
		},
		Deps:    []string{"dep1 v1", "dep2 v1.1"},
		Imports: []string{"github.com/kperreau/goac/scan", "github.com/kperreau/goac/hasher", "anotherlib v1"},
	}
	localDir := ""

	localDeps, extDeps := cleanDeps(rawData, localDir)

	assert.Equal(t, []string{".", "./scan", "./hasher"}, localDeps)
	assert.Equal(t, []string{"dep1 v1", "dep2 v1.1", "anotherlib v1"}, extDeps)
}

func TestCleanDeps_EmptyDepsAndImports(t *testing.T) {
	rawData := &toolData{
		Deps:    []string{},
		Imports: []string{},
	}
	localDir := "../.."

	localDeps, extDeps := cleanDeps(rawData, localDir)

	assert.Equal(t, []string{"../.."}, localDeps)
	assert.Empty(t, extDeps)
}

func TestGetDependencies_ValidModfileAndDependencies(t *testing.T) {
	gomod := &modfile.File{
		Require: []*modfile.Require{
			{Mod: module.Version{Path: "github.com/pkg1", Version: "v1.0.0"}},
			{Mod: module.Version{Path: "github.com/pkg2", Version: "v2.0.0"}},
		},
	}
	rawDeps := []string{"github.com/pkg1", "github.com/pkg2"}
	expected := []string{"github.com/pkg1 v1.0.0", "github.com/pkg2 v2.0.0"}

	deps := getDependencies(gomod, rawDeps)

	if !reflect.DeepEqual(deps, expected) {
		t.Errorf("Expected %v, but got %v", expected, deps)
	}
}

func TestGetDependencies_EmptyDependenciesList(t *testing.T) {
	gomod := &modfile.File{
		Require: []*modfile.Require{
			{Mod: module.Version{Path: "github.com/pkg1", Version: "v1.0.0"}},
			{Mod: module.Version{Path: "github.com/pkg2", Version: "v2.0.0"}},
		},
	}
	var rawDeps []string

	deps := getDependencies(gomod, rawDeps)

	assert.Empty(t, deps)
}

func TestGetDependencies_EmptyModfile(t *testing.T) {
	gomod := &modfile.File{}
	rawDeps := []string{"github.com/pkg1", "github.com/pkg2"}
	expected := rawDeps

	deps := getDependencies(gomod, rawDeps)

	assert.Equal(t, expected, deps)
}

func TestFindVersion_ModuleFound(t *testing.T) {
	dependencies := []*modfile.Require{
		{Mod: module.Version{Path: "module1", Version: "v1.0.0"}},
		{Mod: module.Version{Path: "module2", Version: "v2.0.0"}},
		{Mod: module.Version{Path: "module3", Version: "v3.0.0"}},
	}
	val := "module2"
	expected := "module2 v2.0.0"

	result := findVersion(dependencies, val)

	assert.Equal(t, expected, result)
}

func TestFindVersion_EmptyDependencies(t *testing.T) {
	var dependencies []*modfile.Require
	val := "module1"
	expected := val

	result := findVersion(dependencies, val)

	assert.Equal(t, expected, result)
}

func TestFindVersion_EmptyString(t *testing.T) {
	dependencies := []*modfile.Require{
		{Mod: module.Version{Path: "module1", Version: "v1.0.0"}},
		{Mod: module.Version{Path: "module2", Version: "v2.0.0"}},
		{Mod: module.Version{Path: "module3", Version: "v3.0.0"}},
	}
	val := ""

	result := findVersion(dependencies, val)

	assert.Empty(t, result)
}
