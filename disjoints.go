package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sand8080/go-sketches/db"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func readPassword(c *cli.Context) error {
	if c.GlobalBool("pass") {
		fmt.Println("DB password:")
		pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		c.GlobalSet("password", string(pass))
	}
	return nil
}

func getDBConnection(c *cli.Context) (*sql.DB, error) {
	readPassword(c)
	conn, err := db.GetDBConnection(c.GlobalString("host"),
		c.GlobalInt("port"), c.GlobalString("user"),
		c.GlobalString("password"), c.GlobalString("name"),
		c.GlobalString("ssl"))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func drop(c *cli.Context) error {
	if conn, err := getDBConnection(c); err != nil {
		return err
	} else {
		defer conn.Close()
		return db.DropTables(conn)
	}
}

func create(c *cli.Context) error {
	if conn, err := getDBConnection(c); err != nil {
		return err
	} else {
		defer conn.Close()
		return db.CreateTables(conn)
	}
}

func fill(c *cli.Context) error {
	if conn, err := getDBConnection(c); err != nil {
		return err
	} else {
		defer conn.Close()
		return db.FillTables(conn, c.Int("num"), c.Int("rels"),
			c.Int("min_opr"), c.Int("max_opr"))
	}
}

func recalc(c *cli.Context) error {
	if conn, err := getDBConnection(c); err != nil {
		return err
	} else {
		defer conn.Close()
		return db.RecalculateDisjoints(conn)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "disjoints"
	app.Usage = "Disjoint sets operations manager"
	app.Description = "Tool for create/drop/fill/calculate disjoint sets in DB"
	app.Version = "0.0.1"
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
		{
			Name:   "fill",
			Usage:  "Fill DB",
			Action: fill,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "num, n",
					Usage: "Number of objects in DB",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "rels, r",
					Usage: "Number of relations in DB",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "min_opr, m",
					Usage: "Minimal objects per relation in DB",
					Value: 0,
				},
				cli.IntFlag{
					Name:  "max_opr, M",
					Usage: "Maximal objects per relation in DB",
					Value: 0,
				},
			},
		},
		{
			Name:   "recalc",
			Usage:  "Recalculate disjoints DB",
			Action: recalc,
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
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
		cli.StringFlag{
			Name:  "ssl, s",
			Usage: "SSL connection to DB",
			Value: "disable",
		},
		cli.BoolFlag{
			Name:  "pass, P",
			Usage: "Read DB password from stdin",
		},
		cli.StringFlag{
			Name:   "password",
			Hidden: true,
			Value:  "disjoint",
		},
	}

	app.Action = func(c *cli.Context) error {
		if len(os.Args) < 2 {
			return cli.ShowAppHelp(c)
		}
		readPassword(c)
		return nil
	}
	app.Run(os.Args)
}
