package project

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoadCache_LoadsCacheDataFromFileIfExists(t *testing.T) {
	// Initialize the class object
	p := &Project{
		Cache:    &Cache{},
		HashPath: "hash",
		Path:     "path",
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "goac-cache")
	assert.NoError(t, err)
	DefaultCachePath = tmpDir // override DefaultCachePath

	// Clean up the cache file
	defer os.RemoveAll(tmpDir)

	// Create a temporary cache file
	cacheFilePath := fmt.Sprintf("%s%s.yaml", DefaultCachePath, p.HashPath)
	cacheData := Cache{
		Path: p.Path,
		Target: map[Target]*Metadata{
			TargetBuild: {
				DependenciesHash: "hash",
				DirHash:          "hash",
				Date:             "date",
			},
		},
	}
	data, _ := yaml.Marshal(cacheData)
	err = os.WriteFile(cacheFilePath, data, 0o644)
	assert.NoError(t, err)

	// Call the method under test
	err = p.LoadCache()

	// Assert that the cache data is loaded from file
	assert.NoError(t, err)
	assert.Equal(t, cacheData, *p.Cache)
}

func TestLoadCache_InitializesDefaultCacheIfFileDoesNotExist(t *testing.T) {
	// Initialize the class object
	p := &Project{
		Cache:    &Cache{},
		HashPath: "hash",
		Path:     "path",
	}

	DefaultCachePath = "not-exist" // override DefaultCachePath

	// Call the method under test
	err := p.LoadCache()

	fmt.Println("cache 1:", p.Cache.Target[TargetBuild])

	// Assert that a default cache is initialized
	assert.NoError(t, err)
	assert.Equal(t, p.Path, p.Cache.Path)
	assert.Empty(t, p.Cache.Target)
}

func TestLoadCache_ReturnsErrorIfCacheFileCannotBeRead(t *testing.T) {
	// Initialize the class object
	p := &Project{
		Cache:    &Cache{},
		HashPath: "hash",
		Path:     "path",
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "goac-cache")
	assert.NoError(t, err)
	DefaultCachePath = tmpDir // override DefaultCachePath

	// Clean up the cache file
	defer os.RemoveAll(tmpDir)

	// Create a cache file with no read permissions
	cacheFilePath := fmt.Sprintf("%s%s.yaml", DefaultCachePath, p.HashPath)
	err = os.WriteFile(cacheFilePath, []byte{}, 0o000)
	assert.NoError(t, err)

	// Call the method under test
	err = p.LoadCache()

	// Assert that an error is returned
	assert.Error(t, err)
}

func TestLoadCache_ReturnsErrorIfCacheFileCannotBeUnmarshaled(t *testing.T) {
	// Initialize the class object
	p := &Project{
		Cache:    &Cache{},
		HashPath: "hash",
		Path:     "path",
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "goac-cache")
	assert.NoError(t, err)
	DefaultCachePath = tmpDir // override DefaultCachePath

	// Clean up the cache file
	defer os.RemoveAll(tmpDir)

	// Create a cache file with invalid YAML data
	cacheFilePath := fmt.Sprintf("%s%s.yaml", DefaultCachePath, p.HashPath)
	err = os.WriteFile(cacheFilePath, []byte("invalid_yaml"), 0o644)
	assert.NoError(t, err)

	// Call the method under test
	err = p.LoadCache()

	// Assert that an error is returned
	assert.Error(t, err)
}

func TestReadCacheFromFile_SuccessfulRead(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "goac-cache")
	assert.NoError(t, err)
	DefaultCachePath = tmpDir // override DefaultCachePath

	// Clean up the cache file
	defer os.RemoveAll(tmpDir)

	// Write cache data to the file
	cacheFilePath := fmt.Sprintf("%s%s.yaml", DefaultCachePath, "hash")
	err = os.WriteFile(cacheFilePath, []byte("target:\n    build:\n        dependencieshash: hash1\n        dirhash: hash2\n        date: \"2024-04-15T15:39:47+02:00\""), 0o644)
	assert.NoError(t, err)

	// Create a cache object
	cache := &Cache{
		Target: make(map[Target]*Metadata),
		Path:   cacheFilePath,
	}

	// Call the function under test
	err = readCacheFromFile(cache.Path, cache)
	assert.NoError(t, err)

	// Assert that the cache data was successfully read
	expectedTarget := TargetBuild
	expectedMetadata := &Metadata{
		DependenciesHash: "hash1",
		DirHash:          "hash2",
		Date:             "2024-04-15T15:39:47+02:00",
	}

	assert.NotEmpty(t, cache.Target[expectedTarget])
	assert.Equal(t, expectedMetadata, cache.Target[TargetBuild])
}

func TestReadCacheFromFile_InvalidDataTypes(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "goac-cache")
	assert.NoError(t, err)
	DefaultCachePath = tmpDir // override DefaultCachePath

	// Clean up the cache file
	defer os.RemoveAll(tmpDir)

	// Write cache data to the file
	cacheFilePath := fmt.Sprintf("%s%s.yaml", DefaultCachePath, "hash")
	err = os.WriteFile(cacheFilePath, []byte("target:\n    invalid"), 0o644)
	assert.NoError(t, err)

	// Create a cache object
	cache := &Cache{
		Target: make(map[Target]*Metadata),
		Path:   cacheFilePath,
	}

	// Call the function under test
	err = readCacheFromFile(cache.Path, cache)
	assert.Error(t, err)
}

func TestIsMetadataMatch_ReturnsTrueWhenDependenciesHashAndDirHashMatch(t *testing.T) {
	// Arrange
	cachem := &Metadata{
		DependenciesHash: "hash1",
		DirHash:          "hash2",
	}
	m := &Metadata{
		DependenciesHash: "hash1",
		DirHash:          "hash2",
	}

	// Act
	result := cachem.isMetadataMatch(m)

	// Assert
	assert.True(t, result)
}

func TestIsMetadataMatch_ReturnsFalseWhenDependenciesHashOrDirHashDoNotMatch(t *testing.T) {
	// Arrange
	cachem := &Metadata{
		DependenciesHash: "hash1",
		DirHash:          "hash2",
	}
	m := &Metadata{
		DependenciesHash: "hash3",
		DirHash:          "hash2",
	}

	// Act
	result := cachem.isMetadataMatch(m)

	// Assert
	assert.False(t, result)
}

func TestIsMetadataMatch_PanicWhenBothMetadataObjectsAreNil(t *testing.T) {
	// Arrange
	var cachem *Metadata
	var m *Metadata

	// Assert
	assert.Panics(t, func() { cachem.isMetadataMatch(m) })
}

func TestWriteCache_ValidData(t *testing.T) {
	// Initialize the project object
	p := &Project{
		// Initialize the necessary fields
		Cache: &Cache{
			Target: make(map[Target]*Metadata),
			Path:   DefaultCachePath,
		},
		CMDOptions: &Options{
			Target: TargetBuild,
		},
		Metadata: &Metadata{
			DependenciesHash: "dependenciesHash",
			DirHash:          "dirHash",
		},
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "goac-cache")
	assert.NoError(t, err)
	DefaultCachePath = tmpDir // override DefaultCachePath

	// Clean up the cache file
	defer os.RemoveAll(tmpDir)

	// Call the writeCache method
	err = p.writeCache()

	// Assert that there is no error
	assert.NoError(t, err)

	// Assert that the cache file exists
	_, err = os.Stat(filepath.Join(DefaultCachePath, p.HashPath, ".yaml"))
	assert.NoError(t, err)
}

func TestWriteCache_ErrorCreatingCacheDirectory(t *testing.T) {
	// Initialize the project object
	p := &Project{
		// Initialize the necessary fields
		Cache: &Cache{
			Target: make(map[Target]*Metadata),
			Path:   DefaultCachePath,
		},
		CMDOptions: &Options{
			Target: TargetBuild,
		},
		Metadata: &Metadata{
			DependenciesHash: "dependenciesHash",
			DirHash:          "dirHash",
		},
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "goac-cache")
	assert.NoError(t, err)
	err = os.Mkdir(filepath.Join(tmpDir, "read"), 0o444)
	assert.NoError(t, err)
	DefaultCachePath = filepath.Join(tmpDir, "read") // override DefaultCachePath

	// Clean up the cache file
	defer os.RemoveAll(tmpDir)

	// Call the writeCache method
	err = p.writeCache()

	// Assert that the expected error is returned
	assert.Error(t, err)
}
