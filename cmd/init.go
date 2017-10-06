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
			err error

			noProfile, _ = cmd.Flags().GetBool("no-profile")
			force, _     = cmd.Flags().GetBool("force")
		)

		if err = lanes.InitConfig(noProfile, force); err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
	},
}
