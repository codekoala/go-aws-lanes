package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
	"github.com/codekoala/go-aws-lanes/ssh"
)

var sshCmd = &cobra.Command{
	Use:   "ssh [lane]",
	Short: "List all server names, IP, Instance ID (optionally filtered by lane), prompting for one to SSH into",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		var (
			lane string
			svr  *lanes.Server
			err  error
		)

		if len(args) > 0 {
			lane = args[0]
		}

		if svr, err = ChooseServer(lane); err != nil {
			cmd.Printf(err.Error())
			os.Exit(1)
		}

		if err = ConnectToServer(svr); err != nil {
			cmd.Printf("SSH error: %s\n", err)
			os.Exit(1)
		}
	},
}

// ConnectToServer uses the specified server's lane to correctly connect to the desired server.
func ConnectToServer(svr *lanes.Server, args ...string) (err error) {
	var (
		sshProfile *ssh.Profile
		exists     bool
	)

	if profile == nil {
		return fmt.Errorf("invalid profile selected")
	}

	fmt.Printf("Connecting to server %s...\n", svr)
	if sshProfile, exists = profile.SSH.Mods[svr.Lane]; !exists {
		return fmt.Errorf("No SSH profile for lane %q\n", svr.Lane)
	}

	if err = svr.Login(sshProfile, args); err != nil {
		return fmt.Errorf("connection error: %s\n", err)
	}

	return nil
}
