package main

import (
	"log"
	"os"

	"github.com/programmfabrik/dv/dv"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "dv",
		Usage: "start server for data visualization",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "server",
				Usage: "start visualisation server",
			},
			&cli.StringFlag{
				Name:  "addr",
				Value: ":10000",
				Usage: "set addr for server",
			},
			&cli.StringFlag{
				Name:  "url",
				Value: "http://localhost:10000",
				Usage: "send data url for server",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("server") {
				dv.Server(c)
				// this returns on CTRL-C
			} else {
				dv.Send(c)
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
