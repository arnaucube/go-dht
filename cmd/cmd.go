package cmd

import (
	"fmt"
	"go-dht/config"
	"go-dht/node"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var NodeCommands = []cli.Command{
	{
		Name:   "start",
		Usage:  "start the node",
		Action: cmdStart,
	},
}

func cmdStart(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

	if err := config.MustRead(c); err != nil {
		return err
	}
	fmt.Println(config.C)

	var n node.Node
	var err error
	if config.C.ID != "" {
		n, err = node.LoadNode(config.C.ID)
		if err != nil {
			return err
		}
		log.Info("Node loaded with ID: ", n.ID())
	} else {
		n, err = node.NewNode()
		if err != nil {
			return err
		}
		log.Info("New node created with ID: ", n.ID())
	}

	go func() {
		admin := node.NewAdmin(n)
		err := admin.Start()
		if err != nil {
			panic(err)
		}
	}()

	err = n.Start()
	return err
}
