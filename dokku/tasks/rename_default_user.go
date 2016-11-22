package tasks

import "github.com/strothj/dokku-apply/dokku"

func renameDefaultUserToAdmin() Task {
	return &renameDefaultUserToAdminTask{}
}

type renameDefaultUserToAdminTask struct{}

func (t *renameDefaultUserToAdminTask) Run(config *dokku.Config, env *dokku.Environment) {

}
