package ssh

import (
	"fmt"
	"os"
	"strings"
)

func ReadSSHAuthorizedKey() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	pubSSH := fmt.Sprintf("%s/.ssh/id_rsa.pub", dirname)
	file, err := os.ReadFile(pubSSH)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(file), "\n", ""), nil
}
