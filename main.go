package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/BOBO1997/monkey/repl"
	"github.com/urfave/cli"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "monkey"
	app.Usage = "monkey"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		repl.Start(os.Stdin, os.Stdout)
		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// total 2284 lines
