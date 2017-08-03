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
	fmt.Printf("Current profile: %s\n", Config.Profile)
	if profile, err = Config.GetCurrentProfile(); err != nil {
		fmt.Println(err.Error())
	} else {
		profile.Activate()

		// Create a session to share configuration, and load external configuration.
		sess = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(profile.Region),
		}))

		// Create the service's client with the session.
		svc = ec2.New(sess)
	}

	return RootCmd.Execute()
}
