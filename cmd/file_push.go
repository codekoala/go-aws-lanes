package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
	"github.com/codekoala/go-aws-lanes/ssh"
)

var filePushCmd = &cobra.Command{
	Use:   "push LANE SOURCE... DESTINATION",
	Short: "Pushes each SOURCE file to DESTINATION on all LANE instances",
	Args:  cobra.MinimumNArgs(3),

	Run: func(cmd *cobra.Command, args []string) {
		var (
			servers []*lanes.Server
			err     error
		)

		fl := cmd.Flags()
		args = fl.Args()

		confirmed, _ := cmd.Flags().GetBool("confirm")
		lane := args[0]
		sources := args[1 : len(args)-1]
		dest := args[len(args)-1]

		// make sure the source files exist and are accessible
		if err = CheckSourceFiles(sources...); err != nil {
			cmd.Printf("source file error: %s\n", err)
			os.Exit(1)
		}

		cmd.Printf("Servers that will receive the specified files:\n")
		if servers, err = DisplayLaneAndConfirm(lane, "Type CONFIRM to begin pushing files:", confirmed); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if err = CopyToServers(servers, dest, sources...); err != nil {
			fmt.Printf("copy error: %s\n", err)
			os.Exit(1)
		}
	},
}

func CheckSourceFiles(sources ...string) (err error) {
	for _, src := range sources {
		if _, err = os.Stat(src); err != nil {
			return
		}
	}

	return nil
}

func CopyToServers(servers []*lanes.Server, dest string, sources ...string) (err error) {
	var (
		sshProfile *ssh.Profile
		exists     bool
	)

	if profile == nil {
		return fmt.Errorf("invalid profile selected")
	}

	lane := servers[0].Lane
	if sshProfile, exists = profile.SSH.Mods[lane]; !exists {
		return fmt.Errorf("No SSH profile for lane %q\n", lane)
	}

	for _, svr := range servers {
		fmt.Printf("\n=====\n\nCopying to server %s...\n", svr)

		if err = svr.Push(sshProfile, dest, sources...); err != nil {
			fmt.Printf("connection error: %s\n", err)
			continue
		}
	}

	return nil
}
