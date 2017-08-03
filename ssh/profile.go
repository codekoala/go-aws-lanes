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

func (this *Profile) SSHArgs(addr string) (args []string) {
	if this.User == "" {
		this.User = "ec2-user"
	}

	args = append(args, fmt.Sprintf("%s@%s", this.User, addr))

	if this.Identity != "" {
		args = append(args, "-i", this.Identity)
	}

	if this.Tunnel != "" {
		args = append(args, "-L", this.Tunnel)
	}

	for _, t := range this.Tunnels {
		args = append(args, "-L", t)
	}

	return args
}
