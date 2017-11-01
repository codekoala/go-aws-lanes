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

	PersistentPreRunE: RequireProfile,

	Run: func(cmd *cobra.Command, args []string) {
		var (
			servers []*lanes.Server
			err     error

			fl = cmd.Flags()
		)

		lane := fl.Arg(0)
		shCmd := fl.Arg(1)
		prompt := fmt.Sprintf("\nType CONFIRM to execute %q on these machines:", shCmd)
		confirmed, _ := fl.GetBool("confirm")

		filter, _ := fl.GetString("filter")
		if servers, err = DisplayFilteredLaneAndConfirm(lane, filter, prompt, confirmed); err != nil {
			cmd.Println(err.Error())
			os.Exit(1)
		}

		for _, svr := range servers {
			cmd.Printf("=====\nExecuting on %s (%s):  %s\n", svr.Name, svr.IP, shCmd)
			if err = ConnectToServer(svr, shCmd); err != nil {
				cmd.Printf("SSH error: %s\n", err)
				os.Exit(1)
			}
		}
	},
}
