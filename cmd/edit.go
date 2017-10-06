package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Edit the Lanes profile configuration",
	Args:    cobra.NoArgs,
	Aliases: []string{"ed"},

	PersistentPreRunE: RequireProfile,

	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var editor = os.Getenv("EDITOR")

		if editor == "" {
			editor = "vi"
		}

		cmd.Printf("Editing profile %q using %s\n", Config.Profile, editor)
		ed := exec.Command(editor, lanes.GetProfilePath(Config.Profile, true))
		ed.Stdout = os.Stdout
		ed.Stderr = os.Stderr
		ed.Stdin = os.Stdin

		return ed.Run()
	},
}
