package cmd

import (
	"runtime"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display Lanes version",

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("lanes", version.Version)
		cmd.Println("Commit:\t\t", version.Commit)
		cmd.Println("Build date:\t", version.BuildDate)
		cmd.Println("Go:\t\t", runtime.Version())
	},
}
