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
		cmd.Println("lanes", version.String())
		cmd.Println("Build date:", version.BuildDate)
		cmd.Println("Go:", runtime.Version())
	},
}
