package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/codekoala/go-manidator"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/codekoala/go-aws-lanes"
	"github.com/codekoala/go-aws-lanes/ssh/session"
)

func init() {
	fl := shCmd.Flags()

	fl.BoolP("confirm", "c", false, "Bypass manual confirmation step")
	fl.Bool("parallel", false, "Execute command on each target system in parallel")
	fl.IntP("num-parallel", "n", 0, "Maximum number of target systems to execute command on in parallel")
	fl.Int("pparallel", 0, "Maximum percentage of target systems to execute command on in parallel")
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

		filter, _ := fl.GetString("filter")
		if servers, err = DisplayFilteredLaneAndConfirm(lane, filter, prompt, confirmed); err != nil {
			cmd.Println(err.Error())
			os.Exit(1)
		}

		numParallel := getNumParallel(fl, len(servers))
		cmd.Printf("Executing on %d servers in parallel\n", numParallel)
		if numParallel > 0 {
			result = runInParallel(numParallel, servers, sh)
		} else {
			result = runInSequence(servers, sh)
		}

		if result != nil {
			cmd.Printf("SSH error: %s\n", result)
			os.Exit(1)
		}
	},
}

// getNumParallel determines the number of servers to execute commands on in parallel based on the command line flags.
func getNumParallel(flags *pflag.FlagSet, total int) int {
	inParallel, _ := flags.GetBool("parallel")
	numParallel, _ := flags.GetInt("num-parallel")
	percParallel, _ := flags.GetInt("pparallel")

	if percParallel > 0 {
		// calculate number of instances to hit in parallel
		numParallel = int((float64(percParallel) / 100.0) * float64(total))
	} else if numParallel > 0 {
		// done
	} else if inParallel {
		numParallel = total
	} else {
		numParallel = -1
	}

	// do a little bounds checking
	if numParallel <= 0 {
		numParallel = 1
	} else if numParallel > total {
		numParallel = total
	}

	return numParallel
}

// runInSequence runs the specified command on each server in sequence.
func runInSequence(servers []*lanes.Server, sh string) (result error) {
	for _, svr := range servers {
		fmt.Printf("=====\nExecuting on %s (%s): %s\n", svr.Name, svr.IP, sh)
		if err := ConnectToServer(svr, sh); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

// runInParallel runs the specified command on numParallel servers at the same time.
func runInParallel(numParallel int, servers []*lanes.Server, sh string) (result error) {
	var (
		sessions []*session.Session
		wg       sync.WaitGroup
	)

	ctx, cancel := context.WithCancel(context.Background())

	// used to block commands until there is enough capacity
	limit := make(chan struct{}, numParallel)

	// handle keyboard interrupts
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// immediately add each server to the manidator
	mani := manidator.New()
	for _, svr := range servers {
		sess := session.New(svr)
		sessions = append(sessions, sess)
		mani.Add(sess)
	}

	fmt.Printf("Executing on %d servers: %s\n", len(servers), sh)
	mani.Begin(ctx)

	// launch at most numParallel commands at one time
	for _, s := range sessions {
		limit <- struct{}{}
		wg.Add(1)

		go func(sess *session.Session) {
			if err := sess.Run(ctx, sess.Profile(), sh); err != nil {
				result = multierror.Append(result, err)
			}

			sess.Wait()
			<-limit
			wg.Done()
		}(s)
	}

	// no more commands allowed
	close(limit)

	// wait for either all commands to finish or for lanes to be interruped
	select {
	case <-stop:
		fmt.Println("canceled")
		cancel()
	case <-mani.Done():
		fmt.Println("done")
	}

	// allow everything to finish up cleanly
	mani.Stop()
	wg.Wait()

	// display ALL output from the command each server, grouped by server
	for _, sess := range sessions {
		fmt.Printf("=========== OUTPUT FROM %s ===========\n%s\n", sess.GetName(), sess.String())
	}

	return result
}
