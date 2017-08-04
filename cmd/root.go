package cmd

import (
	"fmt"
	"os"
	"strings"

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
	initCmd.AddCommand(initProfileCmd)

	RootCmd.AddCommand(fileCmd)
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(shCmd)
	RootCmd.AddCommand(sshCmd)
	RootCmd.AddCommand(switchCmd)
	RootCmd.AddCommand(versionCmd)

	filePushCmd.Flags().BoolP("confirm", "c", false, "Bypass manual confirmation step")
	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing configuration")
	initCmd.Flags().BoolP("no-profile", "n", false, "Do not create a default profile")
	shCmd.Flags().BoolP("confirm", "c", false, "Bypass manual confirmation step")
}

func Execute() (err error) {
	isInit := strings.Contains(strings.Join(os.Args, " "), " init")

	if !isInit {
		if Config, err = lanes.LoadConfig(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	return RootCmd.Execute()
}
