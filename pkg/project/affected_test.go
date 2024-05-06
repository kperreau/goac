package project

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestIsAffected_ForceTrue(t *testing.T) {
	p := &Project{
		CMDOptions: &Options{
			Force:  true,
			Target: TargetBuild,
		},
	}
	result := p.isAffected()

	// Assert
	assert.True(t, result)
}

func TestIsAffected_TargetNotInCache(t *testing.T) {
	p := &Project{
		CMDOptions: &Options{
			Target: TargetBuild,
		},
		Cache: &Cache{
			Target: make(map[Target]*Metadata),
		},
	}
	result := p.isAffected()

	// Assert
	assert.True(t, result)
}

func TestIsAffected_BinaryCheckFalse_ReturnFalse(t *testing.T) {
	p := &Project{
		CMDOptions: &Options{
			BinaryCheck: false,
			Target:      TargetBuild,
		},
		Metadata: &Metadata{
			DependenciesHash: "hash",
			DirHash:          "hash",
			Date:             "date",
		},
		CleanPath: ".",
		Name:      "goac",
		Cache: &Cache{
			Target: map[Target]*Metadata{
				TargetBuild: {
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
			},
		},
	}
	result := p.isAffected()

	// Assert
	assert.False(t, result)
}

func TestIsAffected_BinaryCheckTrue_ReturnTrue(t *testing.T) {
	p := &Project{
		CMDOptions: &Options{
			BinaryCheck: true,
			Target:      TargetBuild,
		},
		Metadata: &Metadata{
			DependenciesHash: "hash",
			DirHash:          "hash",
			Date:             "date",
		},
		CleanPath: ".",
		Name:      "goac",
		Cache: &Cache{
			Target: map[Target]*Metadata{
				TargetBuild: {
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
			},
		},
	}
	result := p.isAffected()

	// Assert
	assert.True(t, result)
}

func TestIsAffected_DiffHash_ReturnTrue(t *testing.T) {
	p := &Project{
		CMDOptions: &Options{
			BinaryCheck: true,
			Target:      TargetBuild,
		},
		Metadata: &Metadata{
			DependenciesHash: "new-hash",
			DirHash:          "new-hash",
			Date:             "date",
		},
		CleanPath: ".",
		Name:      "goac",
		Cache: &Cache{
			Target: map[Target]*Metadata{
				TargetBuild: {
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
			},
		},
	}
	result := p.isAffected()

	// Assert
	assert.True(t, result)
}

func TestStringToTarget_Build(t *testing.T) {
	result := StringToTarget(TargetBuild.String())
	assert.Equal(t, TargetBuild, result)
}

func TestStringToTarget_BuildImage(t *testing.T) {
	result := StringToTarget(TargetBuildImage.String())
	assert.Equal(t, TargetBuildImage, result)
}

func TestStringToTarget_EmptyString(t *testing.T) {
	result := StringToTarget("")
	assert.Equal(t, TargetNone, result)
}

func TestStringToTarget_WhitespaceString(t *testing.T) {
	result := StringToTarget("   ")
	assert.Equal(t, TargetNone, result)
}

func TestCountAffected_NoAffectedProjects(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{},
		Options:  &Options{},
	}

	// Invoke the countAffected method
	result := l.countAffected()

	// Assert that the result is 0
	assert.Equal(t, 0, result)
}

func TestCountAffected_OneProject(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{
			{
				CMDOptions: &Options{
					Force:  true,
					Target: TargetBuild,
				},
				Metadata: &Metadata{
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
				CleanPath: ".",
				Name:      "goac",
				Cache: &Cache{
					Target: map[Target]*Metadata{
						TargetBuild: {
							DependenciesHash: "hash",
							DirHash:          "hash",
							Date:             "date",
						},
					},
				},
			},
		},
		Options: &Options{},
	}

	// Invoke the countAffected method
	result := l.countAffected()

	// Assert that the result is 0
	assert.Equal(t, 1, result)
}

func TestCountAffected_OneProjectOfTwo(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{
			{
				CMDOptions: &Options{
					Force:  true,
					Target: TargetBuild,
				},
				Metadata: &Metadata{
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
				CleanPath: ".",
				Name:      "goac",
				Cache: &Cache{
					Target: map[Target]*Metadata{
						TargetBuild: {
							DependenciesHash: "hash",
							DirHash:          "hash",
							Date:             "date",
						},
					},
				},
			},
			{
				CMDOptions: &Options{
					Target: TargetBuild,
				},
				Metadata: &Metadata{
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
				CleanPath: ".",
				Name:      "goac",
				Cache: &Cache{
					Target: map[Target]*Metadata{
						TargetBuild: {
							DependenciesHash: "hash",
							DirHash:          "hash",
							Date:             "date",
						},
					},
				},
			},
		},
		Options: &Options{},
	}

	// Invoke the countAffected method
	result := l.countAffected()

	// Assert that the result is 0
	assert.Equal(t, 1, result)
}

func TestCountAffected_TwoProjects(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{
			{
				CMDOptions: &Options{
					Target: TargetBuild,
				},
				Metadata: &Metadata{
					DependenciesHash: "hash",
					DirHash:          "new-hash",
					Date:             "date",
				},
				CleanPath: ".",
				Name:      "goac",
				Cache: &Cache{
					Target: map[Target]*Metadata{
						TargetBuild: {
							DependenciesHash: "hash",
							DirHash:          "hash",
							Date:             "date",
						},
					},
				},
			},
			{
				CMDOptions: &Options{
					Target: TargetBuild,
				},
				Metadata: &Metadata{
					DependenciesHash: "hash",
					DirHash:          "new-hash",
					Date:             "date",
				},
				CleanPath: ".",
				Name:      "goac",
				Cache: &Cache{
					Target: map[Target]*Metadata{
						TargetBuild: {
							DependenciesHash: "hash",
							DirHash:          "hash",
							Date:             "date",
						},
					},
				},
			},
		},
		Options: &Options{},
	}

	// Invoke the countAffected method
	result := l.countAffected()

	// Assert that the result is 0
	assert.Equal(t, 2, result)
}

func TestAffected_Prints0AffectedProjects(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{},
		Options:  &Options{},
	}

	output, err := redirectAffectedStdout(l.Affected)
	assert.NoError(t, err)

	// Assert that the correct output is printed
	expectedOutput := fmt.Sprintf("Affected: %s/%s\n", color.HiBlueString("%d", 0), color.HiBlueString("%d", len(l.Projects)))
	assert.Equal(t, expectedOutput, output.String())
}

func TestAffected_Prints1AffectedProjects(t *testing.T) {
	// Initialize the List object
	// Initialize the List object
	l := &List{
		Projects: []*Project{
			{
				CMDOptions: &Options{
					Force:          true,
					DryRun:         true,
					Target:         TargetBuild,
					MaxConcurrency: 2,
				},
				Metadata: &Metadata{
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
				CleanPath: ".",
				Name:      "goac",
				Cache: &Cache{
					Target: map[Target]*Metadata{
						TargetBuild: {
							DependenciesHash: "hash",
							DirHash:          "hash",
							Date:             "date",
						},
					},
				},
			},
			{
				CMDOptions: &Options{
					DryRun:         true,
					Target:         TargetBuild,
					MaxConcurrency: 2,
				},
				Metadata: &Metadata{
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
				CleanPath: ".",
				Name:      "goac",
				Cache: &Cache{
					Target: map[Target]*Metadata{
						TargetBuild: {
							DependenciesHash: "hash",
							DirHash:          "hash",
							Date:             "date",
						},
					},
				},
			},
		},
		Options: &Options{
			MaxConcurrency: 2,
			Target:         TargetBuild,
			DryRun:         true,
		},
	}

	output, err := redirectAffectedStdout(l.Affected)
	assert.NoError(t, err)

	// Assert that the correct output is printed
	expectedOutput := fmt.Sprintf("Affected: %s/%s\ngoac => .\n", color.HiBlueString("%d", 1), color.HiBlueString("%d", len(l.Projects)))
	assert.Equal(t, expectedOutput, output.String())
}

func TestAffected_MaxConcurrencyOne(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{
			{
				CMDOptions: &Options{
					Force:          true,
					DryRun:         true,
					Target:         TargetBuild,
					MaxConcurrency: 1,
				},
				Metadata: &Metadata{
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
				CleanPath: ".",
				Name:      "goac",
				Cache: &Cache{
					Target: map[Target]*Metadata{
						TargetBuild: {
							DependenciesHash: "hash",
							DirHash:          "hash",
							Date:             "date",
						},
					},
				},
			},
		},
		Options: &Options{
			MaxConcurrency: 1,
			DryRun:         true,
		},
	}

	// Call the Affected method
	_, err := redirectAffectedStdout(l.Affected)
	assert.NoError(t, err)
}

func TestProcessAffected_BuildAffectedProject(t *testing.T) {
	// Create a project with affected flag set to true and dry run flag set to false
	p := &Project{
		CMDOptions: &Options{
			DryRun:         true,
			Target:         TargetBuild,
			MaxConcurrency: 4,
		},
		Metadata: &Metadata{
			DependenciesHash: "new-hash",
			DirHash:          "hash",
			Date:             "date",
		},
		CleanPath: ".",
		Name:      "goac",
		Cache: &Cache{
			Target: map[Target]*Metadata{
				TargetBuild: {
					DependenciesHash: "hash",
					DirHash:          "hash",
					Date:             "date",
				},
			},
		},
	}

	// Create a processAffectedOptions with a wait group and semaphore
	opts := &processAffectedOptions{
		wg:  &sync.WaitGroup{},
		sem: make(chan bool, 1),
	}

	// Add a wait group counter
	opts.wg.Add(1)

	// Acquire the semaphore
	opts.sem <- true

	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the processAffected function
	processAffected(p, opts)

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

	assert.Equal(t, "goac => .\n", buf.String())
}

func TestPrintAffected_1AffectedAnd1Project(t *testing.T) {
	// Initialize the List object
	l := &List{
		Projects: []*Project{
			{Name: "Project1", CleanPath: "/path/to/project1", CMDOptions: &Options{DryRun: true, Force: true}},
		},
		Options: &Options{
			DryRun: true,
			Force:  true,
		},
	}

	// Call the List method and return output
	buf, _ := redirectAffectedStdout(func() error { l.printAffected(); return nil })

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
				Name: "Project1", CleanPath: "/path/to/project1",
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
	buf, _ := redirectAffectedStdout(func() error { l.printAffected(); return nil })

	// Assert the printed output
	expectedOutput := "Affected: 0/1\n"
	assert.Equal(t, expectedOutput, buf.String())
}

func redirectAffectedStdout(f func() error) (*bytes.Buffer, error) {
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

	return &buf, err
}
