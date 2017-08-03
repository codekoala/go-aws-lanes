package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var (
	RootCmd = &cobra.Command{
		Use:   "lanes",
		Short: "",
		Long:  "",

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
}

func Execute() error {
	fmt.Printf("Current profile: %s\n", Config.Profile)
	profile = Config.GetCurrentProfile()
	profile.Activate()

	// Create a session to share configuration, and load external configuration.
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(profile.Region),
	}))

	// Create the service's client with the session.
	svc = ec2.New(sess)

	return RootCmd.Execute()
}
