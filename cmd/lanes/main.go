package main

import (
	"log"

	"github.com/codekoala/go-aws-lanes"
	"github.com/codekoala/go-aws-lanes/cmd"
)

func main() {
	var err error

	if cmd.Config, err = lanes.LoadConfig(); err != nil {
		log.Fatalln(err.Error())
	}

	if err = cmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
