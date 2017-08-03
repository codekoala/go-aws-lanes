package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var (
	ErrCanceled = errors.New("Canceled")
)

type InputParseFunction func(string) error

func RequireProfile(cmd *cobra.Command, args []string) (err error) {
	fmt.Printf("Current profile: %s\n", Config.Profile)
	if profile, err = Config.GetCurrentProfile(); err != nil || profile == nil {
		fmt.Println("invalid profile selected")
		os.Exit(1)
	}

	profile.Activate()

	// Create a session to share configuration, and load external configuration.
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(profile.Region),
	}))

	// Create the service's client with the session.
	svc = ec2.New(sess)

	return nil
}

func DisplayLaneAndConfirm(lane, prompt string, confirm bool) (servers []*lanes.Server, err error) {
	if servers, err = lanes.FetchServersInLane(svc, lane); err != nil {
		err = fmt.Errorf("failed to fetch servers: %s", err)
		return
	}

	if err = DisplayAndConfirm(servers, prompt, confirm); err != nil {
		return
	}

	return servers, nil
}

func DisplayAndConfirm(servers []*lanes.Server, prompt string, confirm bool) (err error) {
	parse := func(input string) (err error) {
		if input != "CONFIRM" {
			return ErrCanceled
		}

		return nil
	}

	if err = lanes.DisplayServers(servers); err != nil {
		return
	}

	if !confirm {
		if err = Prompt(prompt, parse); err != nil {
			return
		}
	}

	return nil
}

func Prompt(prompt string, parser InputParseFunction) (err error) {
	var input string

	for {
		fmt.Printf("\n%s ", prompt)
		if _, err = fmt.Scanln(&input); err != nil {
			if err == io.EOF {
				goto Cancel
			}

			switch err.Error() {
			case "unexpected newline":
				goto Cancel
			default:
				fmt.Printf("Invalid input: %s\n\n", err)
			}

			continue
		}

		if parser != nil {
			if err = parser(input); err != nil {
				fmt.Printf("Invalid input: %s\n\n", err)
				continue
			}
		}

		break
	}

	return nil

Cancel:
	return ErrCanceled
}

func ChooseServer(lane string) (svr *lanes.Server, err error) {
	var (
		servers []*lanes.Server
		idx     int
	)

	if servers, err = lanes.FetchServersInLane(svc, lane); err != nil {
		err = fmt.Errorf("failed to fetch servers: %s", err)
		return
	}

	if err = lanes.DisplayServers(servers); err != nil {
		err = fmt.Errorf("failed to display servers: %s\n", err)
		return
	}

	parse := func(input string) (err error) {
		if idx, err = strconv.Atoi(input); err != nil {
			return fmt.Errorf("Invalid input; please enter a number.")
		}

		if idx < 1 || idx > len(servers) {
			return fmt.Errorf("Invalid input; please enter a valid server number.")
		}

		return nil
	}

	if err = Prompt("Which server?", parse); err != nil {
		return svr, ErrCanceled
	}

	return servers[idx-1], nil
}
