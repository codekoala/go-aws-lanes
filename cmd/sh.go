package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var shCmd = &cobra.Command{
	Use:   `sh LANE "COMMAND"`,
	Short: "Executes a command on all machines in the specified lane",
	Args:  cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		var (
			lane    string
			shCmd   string
			servers []*lanes.Server
			err     error
		)

		lane = cmd.Flags().Arg(0)
		shCmd = cmd.Flags().Arg(1)

		if servers, err = lanes.FetchServersInLane(svc, lane); err != nil {
			fmt.Printf("failed to fetch servers: %s", err)
			os.Exit(1)
		}

		parse := func(input string) (err error) {
			if input != "CONFIRM" {
				return ErrCanceled
			}

			return nil
		}

		if confirmed, _ := cmd.Flags().GetBool("confirm"); confirmed {
			err = lanes.DisplayServers(servers)
		} else {
			err = Prompt(servers, fmt.Sprintf("Type CONFIRM to execute %q on these machines:", shCmd), parse)
		}

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, svr := range servers {
			fmt.Printf("=====\nExecuting on %s (%s):\t%s\n", svr.Name, svr.IP, shCmd)
			if err = ConnectToServer(svr, shCmd); err != nil {
				fmt.Printf("SSH error: %s\n", err)
				os.Exit(1)
			}
		}
	},
}
