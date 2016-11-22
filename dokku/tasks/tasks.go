package tasks

import "github.com/strothj/dokku-apply/dokku"

// Task is an operation performed as part of configuring Dokku.
type Task interface {
}

func init() {
	dokku.RegisterTask(renameDefaultUserToAdmin())
}
