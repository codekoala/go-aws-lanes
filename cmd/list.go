package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var listCmd = &cobra.Command{
	Use:     "list [lane]",
	Short:   "List all server names, IP, Instance ID (optionally filtered by lane)",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"ls"},

	PersistentPreRunE: RequireProfile,

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			lane    string
			servers []*lanes.Server
		)

		if len(args) > 0 {
			lane = args[0]
		}

		if servers, err = lanes.FetchServersInLane(svc, lane); err != nil {
			return fmt.Errorf("failed to fetch servers: %s\n", err)
		}

		return lanes.DisplayServers(servers)
	},
}
