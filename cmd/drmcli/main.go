package main

import (
	"log"
	"os"

	"github.com/diamondburned/force-broadcast-rgb/cmd/drmcli/drmcli"
)

func main() {
	if err := drmcli.App.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
