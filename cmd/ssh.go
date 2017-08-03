package cmd

import (
	"fmt"
	"os"
	"strconv"

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
			lane    string
			servers []*lanes.Server
			err     error
		)

		if len(args) > 0 {
			lane = args[0]
		}

		if servers, err = lanes.FetchServersInLane(svc, lane); err != nil {
			cmd.Printf("failed to fetch servers: %s", err)
			os.Exit(1)
		}

		svr := ChooseServer(servers)
		if err = ConnectToServer(svr); err != nil {
			cmd.Printf("SSH error: %s\n", err)
			os.Exit(1)
		}
	},
}

func ChooseServer(servers []*lanes.Server) *lanes.Server {
	var (
		idx int
		err error
	)

	if err = lanes.DisplayServers(servers); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	parse := func(input string) (err error) {
		if idx, err = strconv.Atoi(input); err != nil {
			return fmt.Errorf("Invalid input; please enter a number.")
		}

		if idx < 1 || idx > len(servers) {
			return fmt.Errorf("Invalid input; please enter a valid server number.")
		}

		return nil
	}

	if err = Prompt("Which server?", parse); err != nil {
		fmt.Println("Canceled.")
		os.Exit(1)
	}

	return servers[idx-1]
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
