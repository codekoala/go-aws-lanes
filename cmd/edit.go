package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var editCmd = &cobra.Command{
	Use:     "edit [profile]",
	Short:   "Edit the configuration for the current or named Lanes profile",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"ed"},

	PersistentPreRunE: RequireProfile,

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var (
			profile string
			editor  = os.Getenv("EDITOR")
		)

		if len(args) > 0 {
			profile = args[0]
		} else {
			profile = Config.Profile
		}

		if editor == "" {
			editor = "vi"
		}

		cmd.Printf("Editing profile %q using %s\n", profile, editor)
		ed := exec.Command(editor, lanes.GetProfilePath(profile, true))
		ed.Stdout = os.Stdout
		ed.Stderr = os.Stderr
		ed.Stdin = os.Stdin

		return ed.Run()
	},
}
