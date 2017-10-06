package cmd

import (
	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List all lanes profiles",
	Args:  cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		var (
			batch, _ = cmd.Flags().GetBool("batch")
			format   = "  * %s\n"
		)

		if batch {
			format = "%s\n"
		} else {
			cmd.Println("Available profiles:\n")
		}

		for _, profile := range lanes.GetAvailableProfiles() {
			cmd.Printf(format, profile)
		}
	},
}
