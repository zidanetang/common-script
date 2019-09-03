package main

import (
	"github.com/zidanetang/common-script/mongoDB/omogo/handler"
	cli "gopkg.in/urfave/cli.v2"
	"os"
	"time"
)

// VERSION of are
const VERSION = "v0.1.0"

const helpOutput = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}} - {{.Compiled}}
   {{end}}
`

func main() {

	// red := color.New(color.FgRed).SprintFunc()

	cli.AppHelpTemplate = helpOutput

	app := &cli.App{
		Name:     "omogo",
		Usage:    "Insert doucyments into MongoDB",
		Flags:    handler.SetFlags(),
		Compiled: time.Now(),
		Version:  VERSION,
		Action: func(c *cli.Context) error {
			return handler.Run(c)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		handler.PrintError(err)
	}
}
