package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [lane]",
	Short:   "List all server names, IP, Instance ID (optionally filtered by lane)",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"ls"},

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			log.Printf("listing servers filtered by lane %q", args[0])
		} else {
			log.Printf("listing all servers")
		}
	},
}
