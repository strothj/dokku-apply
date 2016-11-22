package dokku

import (
	"os"
	"os/exec"
	"os/user"

	"path"

	"github.com/pkg/errors"
)

// Environment represents the Dokku installation environment.
type Environment struct {
	// Executable is the path to the Dokku executable.
	Executable string

	// UID holds the POSIX Uid for the Dokku system user.
	UID string

	// GID holds the POSIX Gid for the Dokku system user.
	GID string

	// AuthorizedKeys holds the path to the Dokku "authorized_keys". It is
	// normally /home/dokku/.ssh/authorized_keys.
	AuthorizedKeys string

	// AuthorizedKeysMode is the file mode to set on the Dokku authorized_keys
	// file. Set to 0644 by GetEnvironment.
	AuthorizedKeysMode os.FileMode
}

// GetEnvironment returns information about the Dokku installation environment.
func GetEnvironment() (*Environment, error) {
	environment := &Environment{}

	// Test hooks
	lookPath := environmentLookPathTestHook
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	lookup := environmentLookupTestHook
	if lookup == nil {
		lookup = user.Lookup
	}

	executable, err := lookPath("dokku")
	if err != nil {
		return nil, errors.Wrap(err,
			"get environment: failed to find Dokku executable")
	}
	environment.Executable = executable

	user, err := lookup("dokku")
	if err != nil {
		return nil, errors.Wrap(err,
			"get environment: user lookup failed")
	}
	environment.UID = user.Uid
	environment.GID = user.Gid
	environment.AuthorizedKeys = path.Join(user.HomeDir, ".ssh", "authorized_keys")
	environment.AuthorizedKeysMode = os.FileMode(0644)

	return environment, nil
}

var environmentLookPathTestHook func(file string) (string, error)
var environmentLookupTestHook func(username string) (*user.User, error)
