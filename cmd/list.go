package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

func init() {
	listCmd.Flags().BoolP("batch", "b", false, "Batch mode (hide table headers and borders)")
	listCmd.Flags().StringP(
		"columns", "c",
		lanes.GetDefaultColumnList(),
		"Comma-separated list of columns to display")
}

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

		columns, _ := cmd.Flags().GetString("columns")
		parsedColumns := lanes.ParseColumnList(columns)
		if len(parsedColumns) == 0 {
			return fmt.Errorf("invalid columns specified")
		}

		if batch, _ := cmd.Flags().GetBool("batch"); batch {
			config := lanes.GetConfig()
			config.Table.HideTitle = true
			config.Table.HideHeaders = true
			config.Table.HideBorders = true
		}

		if servers, err = lanes.FetchServersInLane(svc, lane); err != nil {
			return fmt.Errorf("failed to fetch servers: %s\n", err)
		}

		return lanes.DisplayServersCols(servers, parsedColumns)
	},
}
