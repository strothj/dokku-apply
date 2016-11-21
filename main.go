package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os/user"

	"strings"

	"log"

	"bytes"

	"github.com/spf13/viper"
)

func main() {
	performUserCheck()
	viper.SetConfigName("dokku-apply")
	viper.AddConfigPath("/etc/")
	viper.ReadInConfig()

	if adminSSHKey := viper.GetString("correctAdminSSHKey"); len(adminSSHKey) > 0 {
		fileContents, err := ioutil.ReadFile("/home/dokku/.ssh/authorized_keys")
		if err != nil {
			log.Fatalln("Unable to read authorized_keys file", err)
		}
		inFile, outFile := bytes.NewBuffer(fileContents), &bytes.Buffer{}
		if err := CorrectAdminSSHKeyName(adminSSHKey, inFile, outFile); err != nil {
			log.Fatalln("Error parsing authorized_keys file", err)
		}
		outFileContents := outFile.String()
		if strings.Compare(string(fileContents), outFileContents) == 0 {
			return
		}
		log.Println("Applying changes to authorized_keys file")
		if err := ioutil.WriteFile("/home/dokku/.ssh/authorized_keys", []byte(outFileContents), 0644); err != nil {
			log.Fatalln("Error committing changes to authorized_keys file", err)
		}
	}
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
