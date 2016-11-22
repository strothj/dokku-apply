package dokku

// Config holds the settings to apply to the Dokku installation.
type Config struct {
	AdminSSHKey string
	Accounts    []struct {
		Name string
		Key  string
	}
}
