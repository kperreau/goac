package project

import (
	"bytes"
	"crypto/sha1"
	"io"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/kperreau/goac/pkg/scan"
	"github.com/stretchr/testify/assert"
)

func TestLoadHashs_Success(t *testing.T) {
	p := &Project{
		Version:  "1.0",
		Name:     "TestProject",
		Path:     ".",
		Target:   make(map[Target]*TargetConfig),
		GoPath:   ".",
		HashPath: ".",
		Module: &Module{
			LocalDirs:    []string{"."},
			ExternalDeps: []string{"dep1", "dep2"},
		},
		HashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
		Metadata:   &Metadata{},
		Cache:      &Cache{},
		Rule:       &scan.Rule{},
		CMDOptions: &Options{},
	}

	err := p.LoadHashs()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	assert.Equal(t, "35380d4ef74486ae75fa80d5f4ba2c3321bf6530", p.Metadata.DependenciesHash)

	if p.Metadata.DirHash == "" {
		t.Error("Expected DirHash to be set, but it was empty")
	}
}

func TestProcessDependenciesHash_ValidProjectWithEmptyExternalDeps(t *testing.T) {
	p := &Project{
		Module: &Module{
			ExternalDeps: []string{},
		},
		HashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
	}

	_, err := processDependenciesHash(p)

	assert.NoError(t, err)
}

func TestProcessDependenciesHash_ValidProjectWithExternalDeps(t *testing.T) {
	p := &Project{
		Module: &Module{
			ExternalDeps: []string{"dep1", "dep2"},
		},
		HashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
	}

	hash, err := processDependenciesHash(p)

	assert.NoError(t, err)
	assert.Equal(t, "35380d4ef74486ae75fa80d5f4ba2c3321bf6530", hash)
}

func TestProcessDirectoryHash_ValidProject_ReturnsHash(t *testing.T) {
	// Arrange
	p := &Project{
		Module: &Module{
			LocalDirs: []string{"../hasher", "../project"},
		},
		Rule: &scan.Rule{},
		CMDOptions: &Options{
			Debug: []string{},
		},
		HashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
	}

	// Act
	result, err := processDirectoryHash(p)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestProcessDirectoryHash_MultipleCallsWithSameProject_ReturnsSameHash(t *testing.T) {
	// Arrange
	p := &Project{
		Module: &Module{
			LocalDirs: []string{"../hasher", "../project"},
		},
		Rule: &scan.Rule{},
		CMDOptions: &Options{
			Debug: []string{},
		},
		HashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
	}

	// Act
	result1, err1 := processDirectoryHash(p)
	result2, err2 := processDirectoryHash(p)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, result1, result2)
}

func TestProcessDirectoryHash_UnableToScanDirectories_ReturnsError(t *testing.T) {
	// Arrange
	p := &Project{
		Module: &Module{
			LocalDirs: []string{"invalid_dir"},
		},
		Rule: &scan.Rule{},
		CMDOptions: &Options{
			Debug: []string{},
		},
		HashPool: &sync.Pool{
			New: func() any { return sha1.New() },
		},
	}

	// Act
	result, err := processDirectoryHash(p)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestDebug_ValidProjectAndFiles_PrintsDebugInformation(t *testing.T) {
	p := &Project{
		Name: "TestProject",
		CMDOptions: &Options{
			Debug: []string{"name", "includes", "excludes", "hashed", "dependencies", "local"},
		},
		Rule: &scan.Rule{
			Includes: []string{"include1", "include2"},
			Excludes: []string{"exclude1", "exclude2"},
		},
		Module: &Module{
			ExternalDeps: []string{"dep1", "dep2"},
			LocalDirs:    []string{"local1", "local2"},
		},
	}
	files := []string{"file1", "file2"}

	output := redirectHashStdout(debug, p, files)

	expectedOutput := "Name: TestProject\n" +
		"Includes\ninclude1\ninclude2\n" +
		"Excludes\nexclude1\nexclude2\n" +
		"Hashed files\nfile1\nfile2\n" +
		"Dependencies\ndep1\ndep2\n" +
		"Local Imports\nlocal1\nlocal2\n\n"

	assert.Equal(t, expectedOutput, output.String())
}

func TestDebug_NoDebugOptionsSpecified_DoesNotPrintDebugInformation(t *testing.T) {
	p := &Project{
		Name: "TestProject",
		CMDOptions: &Options{
			Debug: []string{},
		},
		Rule: &scan.Rule{
			Includes: []string{"include1", "include2"},
			Excludes: []string{"exclude1", "exclude2"},
		},
		Module: &Module{
			ExternalDeps: []string{"dep1", "dep2"},
			LocalDirs:    []string{"local1", "local2"},
		},
	}
	files := []string{"file1", "file2"}

	output := redirectHashStdout(debug, p, files)
	expectedOutput := ""

	assert.Equal(t, expectedOutput, output.String())
}

func TestDebug_EmptyFilesList_PrintsNoHashedFiles(t *testing.T) {
	p := &Project{
		Name: "TestProject",
		CMDOptions: &Options{
			Debug: []string{"hashed"},
		},
		Rule: &scan.Rule{
			Includes: []string{"include1", "include2"},
			Excludes: []string{"exclude1", "exclude2"},
		},
		Module: &Module{
			ExternalDeps: []string{"dep1", "dep2"},
			LocalDirs:    []string{"local1", "local2"},
		},
	}
	var files []string

	output := redirectHashStdout(debug, p, files)
	expectedOutput := "Hashed files\n\n\n"

	assert.Equal(t, expectedOutput, output.String())
}

func redirectHashStdout(f func(*Project, []string), p *Project, files []string) *bytes.Buffer {
	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the func
	f(p, files)

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

	return &buf
}
