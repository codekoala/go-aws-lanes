package cmd

import (
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch PROFILE",
	Short: "Switches AWS profiles (e.g. ~/.lanes/lanes.yml entry)",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		if err := Config.SetProfile(profile); err != nil {
			cmd.Printf("Failed to switch profile: %s\n", err)
		} else {
			cmd.Printf("Switched to profile %q\n", profile)
		}
	},
}
