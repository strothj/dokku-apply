package dokku

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvironment_ExectuableNotFound_ReturnsError(t *testing.T) {
	defer removeEnvironmentTestHooks()
	setEnvironemntTestHooks(t)
	environmentLookPathTestHook = func(file string) (string, error) {
		assert.Equal(t, "dokku", file)
		return "", errors.New("file not found")
	}

	environment, err := GetEnvironment()
	assert.Nil(t, environment)
	assert.Equal(t,
		"get environment: failed to find Dokku executable: file not found",
		err.Error())
}

func TestGetEnvironment_ExecutableFound_ReturnsPath(t *testing.T) {
	defer removeEnvironmentTestHooks()
	setEnvironemntTestHooks(t)

	environment, err := GetEnvironment()
	assert.Nil(t, err)
	assert.Equal(t, "/usr/bin/dokku", environment.Executable)
}

func setEnvironemntTestHooks(t *testing.T) {
	environmentLookPathTestHook = func(file string) (string, error) {
		assert.Equal(t, "dokku", file)
		return "/usr/bin/dokku", nil
	}
}

func removeEnvironmentTestHooks() {
	environmentLookPathTestHook = nil
}
