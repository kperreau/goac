package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugCmd_EmptyStringArg(t *testing.T) {
	arg := ""
	expected := []string{}
	result, err := debugCmd(arg)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestDebugCmd_ValidArgs(t *testing.T) {
	arg := "name,includes,excludes"
	expected := []string{"name", "includes", "excludes"}
	result, err := debugCmd(arg)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestDebugCmd_InvalidArgsOnly(t *testing.T) {
	arg := "invalid1,invalid2"
	expected := []string{}
	errorMsg := "bad debug value: invalid1\nvalid values are: name,includes,excludes,dependencies,local,hashed\n"
	result, err := debugCmd(arg)

	assert.Error(t, err)
	assert.Equal(t, errorMsg, err.Error())
	assert.Equal(t, expected, result)
}

func TestDebugCmd_InvalidArgsWithValid(t *testing.T) {
	arg := "name,invalid1,excludes"
	expected := []string{}
	errorMsg := "bad debug value: invalid1\nvalid values are: name,includes,excludes,dependencies,local,hashed\n"
	result, err := debugCmd(arg)

	assert.Error(t, err)
	assert.Equal(t, errorMsg, err.Error())
	assert.Equal(t, expected, result)
}
