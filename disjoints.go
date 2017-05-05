package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func drop(c *cli.Context) {
	fmt.Println("Deleting tables")
}

func create(c *cli.Context) {
	fmt.Println("Creating tables")
}

func main() {
	app := cli.NewApp()
	app.Name = "disjoints"
	app.Usage = "Disjoint sets operations manager"
	app.Description = "Tool for create/drop/fill/calculate disjoint sets in DB"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host, a",
			Usage: "DB host",
			Value: "localhost",
		},
		cli.IntFlag{
			Name:  "port, p",
			Usage: "DB port",
			Value: 5432,
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "DB name",
			Value: "disjoint",
		},
		cli.StringFlag{
			Name:  "user, u",
			Usage: "DB user",
			Value: "disjoint",
		},
		cli.BoolFlag{
			Name:  "pass, P",
			Usage: "Enter DB password",
		},
		cli.StringFlag{
			Name:   "password",
			Hidden: true,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Create DB",
			Action: create,
		},
		{
			Name:   "drop",
			Usage:  "Drop DB",
			Action: drop,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("pass") {
			pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				panic(err)
			}
			c.Set("password", string(pass))
		}
		return nil
	}
	app.Run(os.Args)
}
