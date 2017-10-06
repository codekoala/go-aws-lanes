package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var autoCompleteCmd = &cobra.Command{
	Use:     "completion [shell]",
	Short:   "Generate shell completion configuration",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"comp"},

	Run: func(cmd *cobra.Command, args []string) {
		var (
			shell = path.Base(os.Getenv("SHELL"))

			fl    = cmd.Flags()
			nargs = fl.NArg()
		)

		if nargs >= 1 {
			shell = fl.Arg(0)
		}

		switch strings.ToLower(shell) {
		case "bash":
			RootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			RootCmd.GenZshCompletion(os.Stdout)
		default:
			cmd.Printf("Unsupported shell: %s\n", shell)
		}
	},
}
