package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCorrectAdminSSHKeyName(t *testing.T) {
	adminSSHKey := "ssh-rsa 40ibbSLwFYImSFGUma9F8FegToN6PGlpm0StsMv6cXIun1yUAbnbINAQe0a3EjrRt4sIDVPPfgNkMtCbMArNAlmy2IjQjAcuUI2OEqPnBu7EIAab7YR7W3DWW7emTefh4gDBpjMQW0SLreD3UX6LH8MzmlexlL8a3SGvjXuELxqINs2yoQpZa6s4Qc5x0TPmXL6sGDw4hzvmPjvxpC7i2h08A4FXLx6CGo9m4E7gYbPAM042BnaNskIotg4NjDv8JrttfybttGYMBhgE6iTIu9qg3BWRJ2qm18R8EynQY145lnaL2nOL8mOQhVFYiC1YpH0ADRnlOOHAJsKRAhF6l2KI3aw8pWaWIIWjLf3WXzpqgcmJcqpa anotherKey"
	originalAuthorizedKeys := `command="FINGERPRINT=SHA256:FoMBitiisGZJQyzDoO5tiItmq4wTh5jsFfPWF5pv0TJ NAME=\"someKey\" ` + "`" + `cat /home/dokku/.sshcommand` + "`" + ` $SSH_ORIGINAL_COMMAND",no-agent-forwarding,no-user-rc,no-X11-forwarding,no-port-forwarding ssh-rsa 40ibbSLwFYImSFGUma9F8FegToN6PGlpm0StsMv6cXIun1yUAbnbINAQe0a3EjrRt4sIDVPPfgNkMtCbMArNAlmy2IjQjAcuUI2OEqPnBu7EIAab7YR7W3DWW7emTefh4gDBpjMQW0SLreD3UX6LH8MzmlexlL8a3SGvjXuELxqINs2yoQpZa6s4Qc5x0TPmXL6sGDw4hzvmPjvxpC7i2h08A4FXLx6CGo9m4E7gYbPAM042BnaNskIotg4NjDv8JrttfybttGYMBhgE6iTIu9qg3BWRJ2qm18R8EynQY145lnaL2nOL8mOQhVFYiC1YpH0ADRnlOOHAJsKRAhF6l2KI3aw8pWaWIIWjLf3WXzpqgcmJcqpa somekey
command="FINGERPRINT=SHA256:FoMBitiisGZJQyzDoO5tiItmq4wTh5jsFfPWF5pv0TJ NAME=\"default\" ` + "`" + `cat /home/dokku/.sshcommand` + "`" + ` $SSH_ORIGINAL_COMMAND",no-agent-forwarding,no-user-rc,no-X11-forwarding,no-port-forwarding ssh-rsa 40ibbSLwFYImSFGUma9F8FegToN6PGlpm0StsMv6cXIun1yUAbnbINAQe0a3EjrRt4sIDVPPfgNkMtCbMArNAlmy2IjQjAcuUI2OEqPnBu7EIAab7YR7W3DWW7emTefh4gDBpjMQW0SLreD3UX6LH8MzmlexlL8a3SGvjXuELxqINs2yoQpZa6s4Qc5x0TPmXL6sGDw4hzvmPjvxpC7i2h08A4FXLx6CGo9m4E7gYbPAM042BnaNskIotg4NjDv8JrttfybttGYMBhgE6iTIu9qg3BWRJ2qm18R8EynQY145lnaL2nOL8mOQhVFYiC1YpH0ADRnlOOHAJsKRAhF6l2KI3aw8pWaWIIWjLf3WXzpqgcmJcqpa anotherKey`
	expectedAuthorizedKeys := `command="FINGERPRINT=SHA256:FoMBitiisGZJQyzDoO5tiItmq4wTh5jsFfPWF5pv0TJ NAME=\"someKey\" ` + "`" + `cat /home/dokku/.sshcommand` + "`" + ` $SSH_ORIGINAL_COMMAND",no-agent-forwarding,no-user-rc,no-X11-forwarding,no-port-forwarding ssh-rsa 40ibbSLwFYImSFGUma9F8FegToN6PGlpm0StsMv6cXIun1yUAbnbINAQe0a3EjrRt4sIDVPPfgNkMtCbMArNAlmy2IjQjAcuUI2OEqPnBu7EIAab7YR7W3DWW7emTefh4gDBpjMQW0SLreD3UX6LH8MzmlexlL8a3SGvjXuELxqINs2yoQpZa6s4Qc5x0TPmXL6sGDw4hzvmPjvxpC7i2h08A4FXLx6CGo9m4E7gYbPAM042BnaNskIotg4NjDv8JrttfybttGYMBhgE6iTIu9qg3BWRJ2qm18R8EynQY145lnaL2nOL8mOQhVFYiC1YpH0ADRnlOOHAJsKRAhF6l2KI3aw8pWaWIIWjLf3WXzpqgcmJcqpa somekey
command="FINGERPRINT=SHA256:FoMBitiisGZJQyzDoO5tiItmq4wTh5jsFfPWF5pv0TJ NAME=\"admin\" ` + "`" + `cat /home/dokku/.sshcommand` + "`" + ` $SSH_ORIGINAL_COMMAND",no-agent-forwarding,no-user-rc,no-X11-forwarding,no-port-forwarding ssh-rsa 40ibbSLwFYImSFGUma9F8FegToN6PGlpm0StsMv6cXIun1yUAbnbINAQe0a3EjrRt4sIDVPPfgNkMtCbMArNAlmy2IjQjAcuUI2OEqPnBu7EIAab7YR7W3DWW7emTefh4gDBpjMQW0SLreD3UX6LH8MzmlexlL8a3SGvjXuELxqINs2yoQpZa6s4Qc5x0TPmXL6sGDw4hzvmPjvxpC7i2h08A4FXLx6CGo9m4E7gYbPAM042BnaNskIotg4NjDv8JrttfybttGYMBhgE6iTIu9qg3BWRJ2qm18R8EynQY145lnaL2nOL8mOQhVFYiC1YpH0ADRnlOOHAJsKRAhF6l2KI3aw8pWaWIIWjLf3WXzpqgcmJcqpa anotherKey`
	inBuff := bytes.NewBuffer([]byte(originalAuthorizedKeys))
	outBuff := &bytes.Buffer{}
	CorrectAdminSSHKeyName(adminSSHKey, inBuff, outBuff)
	if expected, actual := expectedAuthorizedKeys, outBuff.String(); strings.Compare(expected, actual) != 0 {
		t.Fatalf("expected %s\n!=\nactual %s", expected, actual)
	}
}
