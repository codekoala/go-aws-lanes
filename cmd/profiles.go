package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var profilesCmd = &cobra.Command{
	Use:   "profiles [pattern]",
	Short: "List all lanes profiles",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		var (
			fl       = cmd.Flags()
			batch, _ = fl.GetBool("batch")
			format   = "  * %s\n"
			prefix   = ""
		)

		if batch {
			format = "%s\n"
		} else {
			cmd.Println("Available profiles:\n")
		}

		if fl.NArg() > 0 {
			prefix = fl.Arg(0)
		}

		for _, profile := range lanes.GetAvailableProfiles() {
			if strings.HasPrefix(profile, prefix) {
				cmd.Printf(format, profile)
			}
		}

		if !batch {
			cmd.Println()
		}
	},
}
