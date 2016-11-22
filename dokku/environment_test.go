package dokku

import (
	"errors"
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Possibly add validation of returned user.User object.

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

func TestGetEnvironment_UserFound_UID_GID_AuthorizedKeys_AuthorizedKeysMode_Set(t *testing.T) {
	defer removeEnvironmentTestHooks()
	setEnvironemntTestHooks(t)

	environment, err := GetEnvironment()
	assert.Nil(t, err)
	assert.Equal(t, "1000", environment.UID)
	assert.Equal(t, "1000", environment.GID)
	assert.Equal(t, "/home/dokku/.ssh/authorized_keys", environment.AuthorizedKeys)
	assert.Equal(t, os.FileMode(0644), environment.AuthorizedKeysMode)
}

func TestGetEnvironment_UserLookupFailed_ReturnsError(t *testing.T) {
	defer removeEnvironmentTestHooks()
	setEnvironemntTestHooks(t)
	environmentLookupTestHook = func(username string) (*user.User, error) {
		assert.Equal(t, "dokku", username)
		return nil, errors.New("user not found")
	}

	environment, err := GetEnvironment()
	assert.Nil(t, environment)
	assert.Equal(t, "get environment: user lookup failed: user not found", err.Error())
}

func setEnvironemntTestHooks(t *testing.T) {
	environmentLookPathTestHook = func(file string) (string, error) {
		assert.Equal(t, "dokku", file)
		return "/usr/bin/dokku", nil
	}
	environmentLookupTestHook = func(username string) (*user.User, error) {
		assert.Equal(t, "dokku", username)
		return &user.User{
			Uid:      "1000",
			Gid:      "1000",
			Username: "dokku",
			Name:     "dokku",
			HomeDir:  "/home/dokku",
		}, nil
	}
}

func removeEnvironmentTestHooks() {
	environmentLookPathTestHook = nil
	environmentLookupTestHook = nil
}
