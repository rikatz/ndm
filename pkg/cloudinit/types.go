package cloudinit

type Metadata struct {
	InstanceID    string `yaml:"instance-id,omitempty"`
	LocalHostName string `yaml:"local-hostname,omitempty"`
}

type Cmd []string

type Userdata struct {
	Hostname string `yaml:"hostname,omitempty"`
	Users    []User `yaml:"users,omitempty"`
	Runcmd   []Cmd  `yaml:"runcmd,omitempty"`
	Bootcmd  []Cmd  `yaml:"bootcmd,omitempty"`
}

type User struct {
	Name              string   `yaml:"name,omitempty"`
	Groups            string   `yaml:"groups,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
}
