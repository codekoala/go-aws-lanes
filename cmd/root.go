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

		BashCompletionFunction: bashCompletionFunc,

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	Config  *lanes.Config
	profile *lanes.Profile

	sess *session.Session
	svc  *ec2.EC2
)

const bashCompletionFunc = `
__lanes_get_profiles() {
	local out
	out=$(lanes profiles --batch $1 2>&1)
	COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
}

__custom_func() {
	case ${last_command} in
		lanes_switch)
			__lanes_get_profiles
			return
			;;
		*)
			;;
	esac
}
`

func init() {
	annotation := make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string{"__lanes_get_profiles"}

	pfs := RootCmd.PersistentFlags()
	pfs.StringP("profile", "p", "", "use specific profile (for supported commands)")

	profileFlag := pfs.Lookup("profile")
	profileFlag.Annotations = annotation

	RootCmd.AddCommand(autoCompleteCmd)
	RootCmd.AddCommand(editCmd)

	filePushCmd.Flags().BoolP("confirm", "c", false, "Bypass manual confirmation step")
	fileCmd.AddCommand(filePushCmd)
	RootCmd.AddCommand(fileCmd)

	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing configuration")
	initCmd.Flags().BoolP("no-profile", "n", false, "Do not create a default profile")
	initProfileCmd.Flags().BoolP("no-switch", "n", false, "Do not automatically switch to the new profile")
	initCmd.AddCommand(initProfileCmd)

	RootCmd.AddCommand(initCmd)

	listCmd.Flags().StringP("filter", "f", "", "Filter servers by keyword")
	filterFlag := listCmd.Flags().Lookup("filter")
	RootCmd.AddCommand(listCmd)

	profilesCmd.Flags().BoolP("batch", "b", false, "Batch mode")
	RootCmd.AddCommand(profilesCmd)

	shCmd.Flags().AddFlag(filterFlag)
	RootCmd.AddCommand(shCmd)

	sshCmd.Flags().AddFlag(filterFlag)
	RootCmd.AddCommand(sshCmd)

	switchCmd.Annotations = make(map[string]string)
	switchCmd.Annotations[cobra.BashCompCustom] = "__lanes_get_profiles"
	RootCmd.AddCommand(switchCmd)

	RootCmd.AddCommand(versionCmd)
}

func Execute() (err error) {
	isInit := strings.Contains(strings.Join(os.Args, " "), " init")
	isProfileInit := strings.Contains(strings.Join(os.Args, " "), " init profile")

	if !isInit || isProfileInit {
		if err = loadConfig(); err != nil {
			if err != lanes.ErrAbort {
				fmt.Println(err.Error())
			}
			os.Exit(1)
		}
	}

	return RootCmd.Execute()
}

func loadConfig() (err error) {
	if Config, err = lanes.LoadConfig(); err != nil {
		// if we fail to load the configuration, attempt to create it
		if err = lanes.InitConfig(false, false); err != nil {
			return err
		} else {
			// load the newly initialized configuration
			if Config, err = lanes.LoadConfig(); err != nil {
				return err
			}
		}
	}

	return nil
}
