package cloudinit

import (
	"fmt"

	"github.com/rikatz/ndm/pkg/ssh"
	"gopkg.in/yaml.v2"
)

func NewUserdata() Userdata {
	return Userdata{}
}

const (
	headerUserData   = "#cloud-config\n"
	defaultVMMachine = "runner"
)

func NewMetadata() Metadata {
	return Metadata{}
}

func (u *Userdata) PrepareToMachine(hostname, sshkey string) error {
	if u == nil {
		return fmt.Errorf("userdata is null")
	}
	u.Hostname = hostname
	u.Runcmd = []Cmd{
		[]string{"systemctl", "enable", "docker", "--now"},
	}
	user := User{
		Name:   "ndm",
		Groups: "docker",
		SSHAuthorizedKeys: []string{
			sshkey,
		},
	}
	u.Users = append(u.Users, user)

	return nil
}

func (u *Userdata) Render() ([]byte, error) {
	if u == nil {
		return nil, fmt.Errorf("userdata is null")
	}
	yamlValue, err := yaml.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("failure marshalling yaml: %s", err)
	}

	finalBytes := append([]byte(headerUserData), yamlValue...)
	return finalBytes, nil
}

func (m *Metadata) Render() ([]byte, error) {
	if m == nil {
		return nil, fmt.Errorf("userdata is null")
	}
	yamlValue, err := yaml.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failure marshalling yaml: %s", err)
	}
	return yamlValue, nil
}

func GenerateCloudInit() (userdata, metadata []byte, err error) {
	sshKey, err := ssh.ReadSSHAuthorizedKey()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting the SSH Key: %s", err)
	}
	udata := NewUserdata()
	err = udata.PrepareToMachine(defaultVMMachine, sshKey)
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing userdata: %s", err)
	}

	userdata, err = udata.Render()
	if err != nil {
		return nil, nil, fmt.Errorf("error rendering userdata: %s", err)
	}

	mdata := NewMetadata()
	mdata.LocalHostName, mdata.InstanceID = defaultVMMachine, defaultVMMachine
	metadata, err = mdata.Render()
	if err != nil {
		return nil, nil, fmt.Errorf("error rendering metadata: %s", err)
	}
	return userdata, metadata, nil
}
