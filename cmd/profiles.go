package cmd

import (
	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "List all lanes profiles",
	Args:  cobra.NoArgs,

	PersistentPreRunE: RequireProfile,

	Run: func(cmd *cobra.Command, args []string) {
		for _, profile := range lanes.GetAvailableProfiles() {
			cmd.Printf("- %s\n", profile)
		}
	},
}
