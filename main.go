package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os/exec"
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
		if strings.Compare(string(fileContents), outFileContents) != 0 {
			log.Println("Setting admin username to \"admin\" in authorized_keys file")
			if err := ioutil.WriteFile("/home/dokku/.ssh/authorized_keys", []byte(outFileContents), 0644); err != nil {
				log.Fatalln("Error committing changes to authorized_keys file", err)
			}
		}
	}

	if err := UpdateUserList(); err != nil {
		log.Fatalln("Error updating user list", err)
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
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `\"default\"`) {
			if strings.Contains(line, sshKey) {
				line = strings.Replace(line, `\"default\"`, `\"admin\"`, -1)
			}
		}
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

// Account represents a Dokku user account.
type Account struct {
	Name string
	Key  string
}

// UpdateUserList adds missing users and removes left over users. It ignores
// the user account named "admin".
func UpdateUserList() error {
	cmd := exec.Command("dokku", "ssh-keys:list")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	requiredAccounts := make([]Account, 0)
	if err := viper.UnmarshalKey("accounts", &requiredAccounts); err != nil {
		log.Fatalln("Failed to parse user accounts from config file", err)
	}
	currentAccounts := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(output))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		account := scanner.Text()
		if strings.HasPrefix(account, `NAME="`) {
			account = strings.TrimLeft(account, `NAME="`)
			account = strings.TrimRight(account, `"`)
			if account != "admin" {
				currentAccounts = append(currentAccounts, account)
			}
		}
	}
	shouldRemove := func(account string) bool {
		for _, v := range requiredAccounts {
			if account == v.Name {
				return false
			}
		}
		return true
	}
	shouldAdd := func(account string) bool {
		for _, v := range currentAccounts {
			if account == v {
				return false
			}
		}
		return true
	}
	for _, v := range currentAccounts {
		if shouldRemove(v) {
			log.Printf("Removing unneeded account: %s\n", v)
			cmd = exec.Command("dokku", "ssh-keys:remove", v)
			if err := cmd.Run(); err != nil {
				log.Fatalln("Failed to remove unneeded account", v, err)
			}
		}
	}
	for _, v := range requiredAccounts {
		if shouldAdd(v.Name) {
			log.Printf("Adding needed account: %s\n", v.Name)
			cmd = exec.Command("dokku", "ssh-keys:add", v.Name)
			inPipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln("Error opening stdin pipe to ssh-keys:add", err)
			}
			if err := cmd.Start(); err != nil {
				log.Fatalln("Error adding ssh key", err)
			}
			if _, err := inPipe.Write([]byte(v.Key + "\n")); err != nil {
				log.Fatalln("Error piping ssh key to ssh-keys:add")
			}
			if err := cmd.Wait(); err != nil {
				log.Fatalln("Error waiting for ssh-key:add to exit", err)
			}
		}
	}
	return nil
}
