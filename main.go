package main

import (
	"os"

	"go-dht/cmd"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "go-dht"
	app.Version = "0.0.1-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
		cli.BoolFlag{Name: "debug"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, cmd.NodeCommands...)
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
	}
}
