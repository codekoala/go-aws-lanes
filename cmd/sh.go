package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var shCmd = &cobra.Command{
	Use:   "sh [lane]",
	Short: "Executes a command on all machines",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("running command on all servers")
	},
}
