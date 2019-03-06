package main

import (
	"fmt"
	"os"

	"github.com/bukalapak/envsync"
	"github.com/urfave/cli"
)

func main() {
	var source string
	var target string
	syncer := &envsync.Syncer{}

	app := cli.NewApp()
	app.Name = "envsync"
	app.Usage = "synchronize source env and target env file"
	app.UsageText = "envsync -s [source env] -t [target env]"
	app.Version = envsync.VERSION
	app.Copyright = "Bukalapak™ © 2018"
	app.Authors = []cli.Author{
		{
			Name: "PT Bukalapak.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "source, s",
			Usage:       "set source env",
			Value:       "env.sample",
			Destination: &source,
		},
		cli.StringFlag{
			Name:        "target, t",
			Usage:       "set target env",
			Value:       ".env",
			Destination: &target,
		},
	}
	app.Action = func(c *cli.Context) error {
		err := syncer.Sync(source, target)
		if err == nil {
			fmt.Println("source and target are successfully synchronized")
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	app.Run(os.Args)
}
