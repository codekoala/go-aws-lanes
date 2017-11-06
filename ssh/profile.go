package ssh

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
)

var DefaultProfile = Profile{
	User: "ec2-user",
}

type Config struct {
	Default *Profile            `yaml:"default"`
	Mods    map[string]*Profile `yaml:"mods"`
}

type Profile struct {
	User     string   `yaml:"user,omitempty"`
	Identity string   `yaml:"identity,omitempty"`
	Tunnel   string   `yaml:"tunnel,omitempty"`
	Tunnels  []string `yaml:"tunnels,omitempty"`
}

func (this *Profile) GetUser() string {
	if this.User == "" {
		this.User = DefaultProfile.User
	}

	return this.User
}

func (this *Profile) UserAt(addr string) string {
	return fmt.Sprintf("%s@%s", this.GetUser(), addr)
}

func (this *Profile) SSHArgs(addr string) (args []string) {
	if this.Identity != "" {
		this.Identity, _ = homedir.Expand(this.Identity)
		args = append(args, "-i", this.Identity)
	}

	args = append(args, this.UserAt(addr))

	for _, t := range this.AllTunnels() {
		args = append(args, "-L", t)
	}

	return args
}

func (this *Profile) AllTunnels() []string {
	all := this.Tunnels[:]

	if this.Tunnel != "" {
		all = append(all, this.Tunnel)
	}

	return all
}
