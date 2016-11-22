package tasks

import (
	"bufio"
	"io/ioutil"
	"os"

	"bytes"

	"strings"

	"github.com/strothj/dokku-apply/dokku"
)

func renameDefaultUserToAdmin() Task {
	return &renameDefaultUserToAdminTask{}
}

type renameDefaultUserToAdminTask struct{}

func (t *renameDefaultUserToAdminTask) Run(config *dokku.Config, env *dokku.Environment) error {
	readFile := renameDefaultUserReadFileTestHook
	if readFile == nil {
		readFile = ioutil.ReadFile
	}
	writeFile := renameDefaultUserWriteFileTestHook
	if writeFile == nil {
		writeFile = ioutil.WriteFile
	}

	inBytes, err := readFile(env.AuthorizedKeys)
	if err != nil {
		panic("Not Implemented")
	}
	inBuff := bytes.NewBuffer(inBytes)
	outBuff := &bytes.Buffer{}
	scanner := bufio.NewScanner(inBuff)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `\"default\"`) {
			if strings.Contains(line, config.AdminSSHKey) {
				line = strings.Replace(line, `\"default\"`, `\"admin\"`, 1)
			}
		}
		if _, err := outBuff.WriteString(line + "\n"); err != nil {
			panic("Not Implemented")
		}
	}
	if err := scanner.Err(); err != nil {
		panic("Not Implemented")
	}
	if err := writeFile(env.AuthorizedKeys, outBuff.Bytes(), env.AuthorizedKeysMode); err != nil {
		panic("Not Implemented")
	}

	return nil
}

var renameDefaultUserReadFileTestHook func(filename string) ([]byte, error)
var renameDefaultUserWriteFileTestHook func(filename string, data []byte, perm os.FileMode) error
