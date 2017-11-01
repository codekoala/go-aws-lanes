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
	listCmd.Flags().String("hide", "", "Comma-separated list of columns to hide")
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

			fl = cmd.Flags()
		)

		if len(args) > 0 {
			lane = args[0]
		}

		columns, _ := fl.GetString("columns")
		if columns == "" {
			columns = lanes.GetDefaultColumnList()
		}

		parsedColumns := lanes.ParseColumnSet(columns)
		if len(parsedColumns) == 0 {
			return fmt.Errorf("invalid columns specified")
		}

		if hideColumns, _ := fl.GetString("hide"); hideColumns != "" {
			hiddenCols := lanes.ParseColumnSet(hideColumns)
			parsedColumns = parsedColumns.Remove(hiddenCols...)
		}

		if batch, _ := fl.GetBool("batch"); batch {
			lanes.GetConfig().Table.ToggleBatchMode(true)
		}

		filter, _ := fl.GetString("filter")
		if servers, err = lanes.FetchServersInLaneByKeyword(svc, lane, filter); err != nil {
			return fmt.Errorf("failed to fetch servers: %s\n", err)
		}

		return lanes.DisplayServersCols(servers, parsedColumns)
	},
}
