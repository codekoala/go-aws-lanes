package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var initProfileCmd = &cobra.Command{
	Use:   "profile [NAME] [AWS ACCESS KEY ID]",
	Short: "Initialize a new lane profile called NAME",
	Args:  cobra.MaximumNArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		var (
			name string
			err  error

			fl      = cmd.Flags()
			nargs   = fl.NArg()
			profile = lanes.GetSampleProfile()
		)

		if nargs >= 1 {
			name = fl.Arg(0)
		}

		if nargs == 2 {
			profile.AWSAccessKeyId = fl.Arg(1)
		}

		// prompt for profile name
		if name == "" {
			parseName := func(input string) error {
				name = input
				return nil
			}

			if err = Prompt("Profile name?", parseName); err != nil {
				cmd.Printf("Error: %s\n", err)
				os.Exit(1)
			}
		}

		// prompt for AWS access key
		if profile.AWSAccessKeyId == "" {
			parseAccessKey := func(input string) error {
				profile.AWSAccessKeyId = input
				return nil
			}

			if err = Prompt("AWS access key ID?", parseAccessKey); err != nil {
				cmd.Printf("Error: %s\n", err)
				os.Exit(1)
			}
		}

		// prompt for AWS secret key
		parseSecretKey := func(input string) error {
			profile.AWSSecretAccessKey = input
			return nil
		}

		if err = PromptHideInput("AWS secret access key?", parseSecretKey, true); err != nil {
			cmd.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		cmd.Printf("Creating new profile: %s\n", name)
		if err = profile.Write(name); err != nil {
			cmd.Printf("Failed to write default profile: %s\n", err)
			os.Exit(1)
		}
	},
}
