package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch PROFILE",
	Short: "Switches AWS profiles (e.g. ~/.lanes/lanes.yml entry)",
	Args:  cobra.ExactArgs(1),
	RunE:  SwitchToProfile,
}

func SwitchToProfile(cmd *cobra.Command, args []string) (err error) {
	profile := args[0]

	if err = Config.SetProfile(profile); err != nil {
		cmd.Printf("Failed to switch profile: %s\n", err)
		os.Exit(1)
	} else {
		cmd.Printf("Switched to profile %q\n", profile)
	}

	return nil
}
