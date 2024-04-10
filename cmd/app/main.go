package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vuquang23/poseidon/cmd/app/api"
)

func main() {
	app := &cli.App{
		Name: "Poseidon",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "internal/pkg/config/default.yaml",
				Usage:   "Configuration file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "API Server",
				Aliases: []string{"api"},
				Usage:   "Run API Server",
				Action:  api.RunAPI,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
