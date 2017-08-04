package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Lanes",
	Args:  cobra.ExactArgs(0),

	Run: func(cmd *cobra.Command, args []string) {
		var (
			cfg *lanes.Config
			err error

			noProfile, _ = cmd.Flags().GetBool("no-profile")
			force, _     = cmd.Flags().GetBool("force")
		)

		if _, err = os.Stat(lanes.CONFIG); err == nil {
			cmd.Printf("Lanes already appears to be configured! ")
			if !force {
				cmd.Println("Aborting.")
				os.Exit(1)
			} else {
				cmd.Println("Overwriting existing configuration.")
			}
		}

		if cfg, err = lanes.LoadConfigBytes([]byte("profile: default")); err != nil {
			cmd.Printf("Failed to initialize configuration: %s\n", err)
			os.Exit(1)
		}

		if err = cfg.Write(); err != nil {
			cmd.Printf("Failed to write configuration: %s\n", err)
			os.Exit(1)
		}

		if !noProfile {
			p := lanes.GetSampleProfile()
			if err = p.Write("default"); err != nil {
				cmd.Printf("Failed to write default profile: %s\n", err)
				os.Exit(1)
			}
		}
	},
}
