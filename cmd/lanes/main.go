package main

import (
	"fmt"
	"os"

	"github.com/codekoala/go-aws-lanes/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
