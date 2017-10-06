package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/codekoala/go-aws-lanes"
)

var (
	ErrCanceled       = errors.New("Canceled")
	ErrInvalidProfile = errors.New("invalid profile selected")
)

// InputParseFunction deal with validating and saving user input to the appropriate variable.
type InputParseFunction func(string) error

// RequireProfile ensures that a valid profile is configured before allowing certain commands to proceed.
func RequireProfile(cmd *cobra.Command, args []string) (err error) {
	fmt.Printf("Current profile: %s\n", Config.Profile)
	if profile, err = Config.GetCurrentProfile(); err != nil || profile == nil {
		fmt.Println(ErrInvalidProfile)
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

// DisplayLaneAndConfirm displays a table of all instances in the specified lane and requires the user to confirm their
// intentions before allowing the calling operation to continue.
func DisplayLaneAndConfirm(lane, prompt string, confirm bool) (servers []*lanes.Server, err error) {
	if servers, err = profile.FetchServersInLane(svc, lane); err != nil {
		err = fmt.Errorf("failed to fetch servers: %s", err)
		return
	}

	if err = DisplayAndConfirm(servers, prompt, confirm); err != nil {
		return
	}

	return servers, nil
}

// DisplayAndConfirm displays a table with the specified servers and requires the user to confirm their intentions
// before allowing the calling operation to continue.
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

// ReadInput accepts user input, optionally hiding user input for sensitive data.
func ReadInput(hideInput bool) (out string, err error) {
	if hideInput {
		var pw []byte
		if pw, err = terminal.ReadPassword(int(syscall.Stdin)); err != nil {
			return
		}

		fmt.Printf("\n")
		out = string(pw)
	} else {
		if _, err = fmt.Scanln(&out); err != nil {
			return
		}
	}

	return out, nil
}

// Prompt displays a regular prompt for user input.
func Prompt(prompt string, parser InputParseFunction) (err error) {
	return PromptHideInput(prompt, parser, false)
}

// PromptHideInput displays a prompt for user input, optionally hiding user input. The prompt is displayed until the
// user's input is deemed valid.
func PromptHideInput(prompt string, parser InputParseFunction, hideInput bool) (err error) {
	var input string

	for {
		fmt.Printf("%s ", prompt)
		if input, err = ReadInput(hideInput); err != nil {
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

// ChooseServer displays a table of all instances in the specified lane and prompts the user to select one server
// before proceeding with the calling operation.
func ChooseServer(lane string) (svr *lanes.Server, err error) {
	var (
		servers   []*lanes.Server
		selection *lanes.Server
		idx       int
	)

	if servers, err = profile.FetchServersInLane(svc, lane); err != nil {
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

	switch len(servers) {
	case 0:
		return nil, ErrCanceled
	case 1:
		selection = servers[0]
	default:
		prompt := fmt.Sprintf("\nWhich server (1-%d)?", len(servers))
		if err = Prompt(prompt, parse); err != nil {
			return svr, ErrCanceled
		}

		selection = servers[idx-1]
	}

	return selection, nil
}
