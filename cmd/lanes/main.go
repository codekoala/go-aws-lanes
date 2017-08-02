package main

import (
	"log"
	"os"

	"github.com/codekoala/go-aws-lanes"
)

var config *lanes.Config

func init() {
	if len(os.Args) < 2 {
		log.Fatalln("usage: lanes [lane]")
	}
}

func main() {
	var err error

	if config, err = lanes.LoadConfig(); err != nil {
		log.Fatalln(err.Error())
	}

	lane := os.Args[1]
	prof := config.GetCurrentProfile()
	log.Printf("SSH comand: %#v", prof.SSH.Mods[lane].GetCommand("foo"))
}
