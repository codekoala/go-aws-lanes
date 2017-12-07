package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/codekoala/go-manidator"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/codekoala/go-aws-lanes"
	"github.com/codekoala/go-aws-lanes/ssh/session"
)

func init() {
	fl := shCmd.Flags()

	fl.BoolP("confirm", "c", false, "Bypass manual confirmation step")
	fl.Bool("parallel", false, "Execute command on each target system in parallel")
}

var shCmd = &cobra.Command{
	Use:   `sh LANE "COMMAND"`,
	Short: "Executes a command on all machines in the specified lane",
	Args:  cobra.ExactArgs(2),

	PersistentPreRunE: RequireProfile,

	Run: func(cmd *cobra.Command, args []string) {
		var (
			servers []*lanes.Server
			err     error
			result  error

			fl = cmd.Flags()
		)

		lane := fl.Arg(0)
		sh := fl.Arg(1)
		prompt := fmt.Sprintf("\nType CONFIRM to execute %q on these machines:", sh)
		confirmed, _ := fl.GetBool("confirm")
		inParallel, _ := fl.GetBool("parallel")

		filter, _ := fl.GetString("filter")
		if servers, err = DisplayFilteredLaneAndConfirm(lane, filter, prompt, confirmed); err != nil {
			cmd.Println(err.Error())
			os.Exit(1)
		}

		if inParallel {
			result = runInParallel(servers, sh)
		} else {
			result = runInSequence(servers, sh)
		}

		if result != nil {
			cmd.Printf("SSH error: %s\n", result)
			os.Exit(1)
		}
	},
}

func runInSequence(servers []*lanes.Server, sh string) (result error) {
	for _, svr := range servers {
		fmt.Printf("=====\nExecuting on %s (%s): %s\n", svr.Name, svr.IP, sh)
		if err := ConnectToServer(svr, sh); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

func runInParallel(servers []*lanes.Server, sh string) (result error) {
	var sessions []*session.Session

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())

	fmt.Printf("Executing on %d servers: %s\n", len(servers), sh)
	mani := manidator.New()
	for _, svr := range servers {
		sess := session.New(svr)
		if err := sess.Run(ctx, sh); err != nil {
			result = multierror.Append(result, err)
			continue
		}

		mani.Add(sess)
		sessions = append(sessions, sess)
	}

	mani.Begin(ctx)
	select {
	case <-stop:
		fmt.Println("canceled")
		cancel()
	case <-mani.Done():
		fmt.Println("done")
	}
	mani.Stop()

	for _, sess := range sessions {
		fmt.Printf("=========== OUTPUT FROM %s ===========\n%s\n", sess.GetName(), sess.String())
	}

	return result
}
