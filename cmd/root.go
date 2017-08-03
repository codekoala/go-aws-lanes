package cmd

import (
	"log"
	"strings"

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
	profile = Config.GetCurrentProfile()
	profile.Activate()

	// Create a session to share configuration, and load external configuration.
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(profile.Region),
	}))

	// Create the service's client with the session.
	svc = ec2.New(sess)

	out, err := FetchServersInLane("prod")
	log.Printf("%#v %s", err, out)

	return RootCmd.Execute()
}

func FetchServers() ([]*lanes.Server, error) {
	return FetchServersBy(nil)
}

func FetchServersInLane(lane string) ([]*lanes.Server, error) {
	return FetchServersBy(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{{
			Name:   aws.String("tag-key"),
			Values: []*string{aws.String("Lane")},
		}, {
			Name:   aws.String("tag-value"),
			Values: []*string{aws.String(lane)},
		}},
	})
}

func FetchServersBy(input *ec2.DescribeInstancesInput) (servers []*lanes.Server, err error) {
	var out *ec2.DescribeInstancesOutput

	if out, err = svc.DescribeInstances(input); err != nil {
		return
	}

	for _, rez := range out.Reservations {
		for _, inst := range rez.Instances {
			if inst.PublicIpAddress == nil || *inst.PublicIpAddress == "" {
				continue
			}

			svr := &lanes.Server{
				ID: *inst.InstanceId,
				IP: *inst.PublicIpAddress,
			}

			for _, tag := range inst.Tags {
				switch strings.ToLower(*tag.Key) {
				case "name":
					svr.Name = *tag.Value
				case "lane":
					svr.Lane = *tag.Value
				}
			}

			servers = append(servers, svr)
		}
	}

	return servers, nil
}
