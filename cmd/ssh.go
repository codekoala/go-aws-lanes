package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
)

var sshCmd = &cobra.Command{
	Use:   "ssh [lane]",
	Short: "List all server names, IP, Instance ID (optionally filtered by lane), prompting for one to SSH into",
	Args:  cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		var (
			lane    string
			servers []*lanes.Server
			input   string
			idx     int
			err     error
		)

		if len(args) > 0 {
			lane = args[0]
		}

		if servers, err = lanes.FetchServersInLane(svc, lane); err != nil {
			cmd.Printf("failed to fetch servers: %s", err)
			os.Exit(1)
		}

		lanes.DisplayServers(servers)

		for {
			fmt.Printf("\nWhich server? ")
			if _, err = fmt.Scanln(&input); err != nil {
				switch err.Error() {
				case "unexpected newline":
					cmd.Println("Canceled.")
					os.Exit(0)
				default:
					cmd.Printf("Invalid input: %s\n\n", err)
				}

				continue
			}

			if idx, err = strconv.Atoi(input); err != nil {
				cmd.Println("Invalid input; please enter a number.")
				continue
			}

			if idx < 1 || idx > len(servers) {
				cmd.Println("Invalid input; please enter a valid server number.")
				continue
			}

			break
		}

		fmt.Println("CONNECTING TO SERVER", idx)
	},
}
