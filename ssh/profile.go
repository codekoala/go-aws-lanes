package ssh

import (
	"fmt"
)

type Config struct {
	Mods map[string]*Profile `yaml:"mods"`
}

type Profile struct {
	User     string   `yaml:"user,omitempty"`
	Identity string   `yaml:"identity,omitempty"`
	Tunnel   string   `yaml:"tunnel,omitempty"`
	Tunnels  []string `yaml:"tunnels,omitempty"`
}

func (this *Profile) GetCommand(addr string) (cmd string) {
	if this.User == "" {
		this.User = "ec2-user"
	}

	cmd = fmt.Sprintf("ssh %s@%s", this.User, addr)

	if this.Identity != "" {
		cmd += fmt.Sprintf(" -i %s", this.Identity)
	}

	if this.Tunnel != "" {
		cmd += fmt.Sprintf(" -L %s", this.Tunnel)
	}

	for _, t := range this.Tunnels {
		cmd += fmt.Sprintf(" -L %s", t)
	}

	return cmd
}
