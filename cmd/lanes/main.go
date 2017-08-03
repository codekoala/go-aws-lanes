package main

import (
	"fmt"
	"os"

	"github.com/codekoala/go-aws-lanes"
	"github.com/codekoala/go-aws-lanes/cmd"
)

func main() {
	var err error

	if cmd.Config, err = lanes.LoadConfig(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
