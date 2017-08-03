package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var (
	RootCmd = &cobra.Command{
		Use:   "lanes",
		Short: "Helper for interacting with sets of AWS EC2 instances",

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	Config  *lanes.Config
	profile *lanes.Profile

	sess *session.Session
	svc  *ec2.EC2
)

func init() {
	fileCmd.AddCommand(filePushCmd)

	RootCmd.AddCommand(fileCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(shCmd)
	RootCmd.AddCommand(sshCmd)
	RootCmd.AddCommand(switchCmd)

	shCmd.Flags().BoolP("confirm", "c", false, "Bypass manual confirmation step")
	filePushCmd.Flags().BoolP("confirm", "c", false, "Bypass manual confirmation step")
}

func Execute() (err error) {
	return RootCmd.Execute()
}
