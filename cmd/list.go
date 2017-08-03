package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var listCmd = &cobra.Command{
	Use:     "list [lane]",
	Short:   "List all server names, IP, Instance ID (optionally filtered by lane)",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"ls"},

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
			cmd.Printf("failed to fetch servers: %s\n", err)
			os.Exit(1)
		}

		lanes.DisplayServers(servers)
	},
}
