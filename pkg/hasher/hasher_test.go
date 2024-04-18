package hasher

import (
	"crypto/sha1"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiles_HashOfAllFilesContentsInAlphabeticalOrder(t *testing.T) {
	files, err := listFiles("../project")
	assert.NoError(t, err)

	// Call the Files function
	hashPool := NewPool()
	result, err := Files(files, hashPool)

	// Verify the result
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// Cal 2 times to be sure that the result is identical
	result2, err2 := Files(files, hashPool)

	// Verify the result
	assert.NoError(t, err2)
	assert.NotEmpty(t, result2)
}

func listFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, filepath.Join(dir, file.Name()))
		}
	}
	return fileNames, nil
}

func TestWithPool_ReturnsHashOfInputString(t *testing.T) {
	// Create a new hashPool
	hashPool := NewPool()

	// Create the input string
	input := "test string"

	// Call the function under test
	result, err := WithPool(hashPool, input)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the result matches the expected hash
	assert.Equal(t, "661295c9cbf9d6b2f6428414504a8deed3020641", result)
}

func TestNewPool_ReturnsSyncPoolWithSha1New(t *testing.T) {
	pool := NewPool()

	assert.NotNil(t, pool)
	assert.NotNil(t, pool.New)
	assert.IsType(t, sha1.New(), pool.New())
}
