package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/codekoala/go-aws-lanes"
)

var (
	ErrCanceled = errors.New("Canceled")
)

type InputParseFunction func(string) error

func Prompt(servers []*lanes.Server, prompt string, parser InputParseFunction) (err error) {
	var input string

	if err = lanes.DisplayServers(servers); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

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
