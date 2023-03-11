package main

import (
	"os"

	"github.com/urfave/cli/v2"

	_ "github.com/forbot161602/pbc-golang-lib/source/entry/gbpreset"
	"github.com/forbot161602/pbc-golang-lib/source/entry/gbshow_info"
)

var (
	app *cli.App
)

func init() {
	app = &cli.App{
		Name:      "golang-lib",
		Usage:     "Golang library",
		Version:   "v1",
		HelpName:  "./main.exe",
		ArgsUsage: "[arguments...]",
		Authors: []*cli.Author{
			&cli.Author{Name: "forbot161602@gmail.com"},
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:      "show-info",
				Usage:     "Present service information",
				HelpName:  "show-info",
				ArgsUsage: "[arguments...]",
				Action: func(c *cli.Context) error {
					return gbshow_info.Execute()
				},
			},
		},
	}
}

func main() {
	app.Run(os.Args)
}
