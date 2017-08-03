package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch PROFILE",
	Short: "Switches AWS profiles (e.g. ~/.lanes/lanes.yml entry)",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]
		Config.SetProfile(profile)
		if err := Config.Write(); err != nil {
			log.Printf("Failed to switch profile: %s", err)
		} else {
			log.Printf("Switched to profile %q", profile)
		}
	},
}
