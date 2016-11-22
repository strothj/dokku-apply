package tasks

import (
	"strings"
	"testing"

	"os"

	"reflect"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	"github.com/strothj/dokku-apply/dokku"
)

func TestRenameDefaultUser_UpdatesFileContents(t *testing.T) {
	adminKey, startingContents, expectedContents := createRenameDefaultUsersTestValues()
	config := &dokku.Config{AdminSSHKey: adminKey}
	env := &dokku.Environment{
		UID:                "1000",
		GID:                "1000",
		AuthorizedKeys:     "/home/dokku/.ssh/authorized_keys",
		AuthorizedKeysMode: os.FileMode(0644),
	}

	defer removeRenameDefaultUserTestHooks()
	setRenameDefaultUserReadFileTestHook(t, startingContents)
	var writeFileCalled bool
	renameDefaultUserWriteFileTestHook = func(filename string, data []byte, perm os.FileMode) error {
		writeFileCalled = true
		assert.Equal(t, "/home/dokku/.ssh/authorized_keys", filename)
		assert.True(t, reflect.DeepEqual([]byte(expectedContents), data))
		assert.Equal(t, os.FileMode(0644), perm)
		return nil
	}

	task := &renameDefaultUserToAdminTask{}
	err := task.Run(config, env)
	assert.Nil(t, err)
	assert.True(t, writeFileCalled, "WriteFile should have been called")
}

func setRenameDefaultUserReadFileTestHook(t *testing.T, startingContent string) {
	renameDefaultUserReadFileTestHook = func(filename string) ([]byte, error) {
		assert.Equal(t, "/home/dokku/.ssh/authorized_keys", filename)
		return []byte(startingContent), nil
	}
}

func removeRenameDefaultUserTestHooks() {
	renameDefaultUserReadFileTestHook = nil
}

// createRenameDefaultUsersTestValues returns the SSH key for the admin (key),
// the contents of the authorized_keys before it is correct (before), and the
// contents afterwards (after).
func createRenameDefaultUsersTestValues() (key, before, after string) {
	// Generate a line that would appear in the authorized_keys file:
	// command="FINGERPRINT=SHA256:GMWJGSAPGATLMODYLUMGUQDIWWMNPPQIMGEGMGKBNPT
	// NAME=\"test\" `cat /home/dokku/.sshcommand` $SSH_ORIGINAL_COMMAND",
	// no-agent-forwarding,no-user-rc,no-X11-forwarding,no-port-forwarding
	// ssh-rsa NIWDYJDHVOYAFNNDFBKFGQDLCIVNBCSWTCNVUXDMIQJISQCDLVHPDCSFXABPHMKKK
	// PHXDQHJWRMVTBCGEQUSLTGCAKKMYDSJDYWMHQOPLBBPBTOFIBPCIRGDEWXRISSQKKPEMHPVOF
	// NAJUIUJWKPMEUXGUQAGERGCLUYMWRNLDNSGTHHCAGDURSUVLHUKWELKUEEOKBGHQAOMILEKJP
	// TOPMPAJXKSVPYTULNWTUYMYSGXCUWVRBPHEREGVFLSDQJRIIDDYJAYGMLROIJVXYMHYOGVIDQ
	// BCMLBFMOXMFJDJLCSUVDPUGFBJQGUFJCRDASFIGPIJRXMXWDICVXNGYGPGHCCEIFLMMEVSXGP
	// WUTGYYUTJHYHDYA key_comment
	generateLine := func(username string) (line, key string) {
		line = `command="FINGERPRINT=SHA256:`
		line += randomdata.Letters(43) // Gibberish fingerprint
		line += " "
		line += `NAME=\"` + username + `\"`
		line += " `"
		line += "cat /home/dokku/.sshcommand"
		line += "` "
		line += `$SSH_ORIGINAL_COMMAND",no-agent-forwarding,no-user-rc,no-X11-forwarding,no-port-forwarding ssh-rsa `
		key = randomdata.Letters(372)
		line += key // Gibberish RSA public key
		line += " key_comment\n"
		key = "ssh-rsa " + key + " key_comment"
		return
	}

	line, _ := generateLine(randomdata.SillyName())
	before += line
	after += line
	line, key = generateLine("default")
	before += line
	after += strings.Replace(line, "default", "admin", 1)
	return
}
