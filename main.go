package main

import (
	"bufio"
	"io"
	"os/user"

	"strings"

	"github.com/spf13/viper"
)

func main() {
	performUserCheck()
	viper.SetConfigName("dokku-apply")
	viper.ReadInConfig()
}

func performUserCheck() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	if user.Uid != "0" {
		panic("Command must be run as root.")
	}
}

// CorrectAdminSSHKeyName sets the name of the SSH key for the administrator to
// "admin" if it was set to "default". It transforms only the first occurance.
func CorrectAdminSSHKeyName(sshKey string, inFile io.Reader, outFile io.Writer) error {
	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)
	lineEnding := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `\"default\"`) {
			if strings.Contains(line, sshKey) {
				line = strings.Replace(line, `\"default\"`, `\"admin\"`, -1)
			}
		}
		if _, err := writer.WriteString(lineEnding + line); err != nil {
			return err
		}
		lineEnding = "\n"
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}
