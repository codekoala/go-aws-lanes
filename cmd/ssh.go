package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var sshCmd = &cobra.Command{
	Use:   "ssh [lane]",
	Short: "List all server names, IP, Instance ID (optionally filtered by lane), prompting for one to SSH into",
	Args:  cobra.MaximumNArgs(1),

	PersistentPreRunE: RequireProfile,

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
			cmd.Println(err.Error())
			os.Exit(1)
		}

		if err = ConnectToServer(svr); err != nil {
			cmd.Println(err.Error())
			os.Exit(1)
		}
	},
}

// ConnectToServer uses the specified server's lane to correctly connect to the desired server.
func ConnectToServer(svr *lanes.Server, args ...string) (err error) {
	fmt.Printf("Connecting to server %s...\n", svr)
	if err = svr.Login(args); err != nil {
		return fmt.Errorf("connection error: %s", err)
	}

	return nil
}
