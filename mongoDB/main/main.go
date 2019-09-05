package main

import (
	"context"
	"fmt"
	"github.com/zidanetang/common-script/mongoDB/omogo/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	cli "gopkg.in/urfave/cli.v2"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func SetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "servers",
			Aliases: []string{"s"},
			Usage: "Specify MongoDB server Cluster, example, e.g.\n" +
				"\t\t\t -s 127.0.0.1:27017\n" +
				"\t\t\t -s 127.0.0.1:27017,127.0.0.1:27016,127.0.0.1:27015",
		},
		&cli.StringFlag{
			Name:    "counts",
			Aliases: []string{"c"},
			Usage:   "Specify insert documents counts",
		},
		&cli.StringFlag{
			Name:    "DB",
			Aliases: []string{"d"},
			Usage:   "Database name",
		},
		&cli.StringFlag{
			Name:    "Collection",
			Aliases: []string{"t"},
			Usage:   "Collection name",
		},
	}
}

func Run(c *cli.Context) error {
	servers := c.String("servers")
	count := c.String("counts")
	db := c.String("DB")
	coll := c.String("Collection")

	uri := "mongodb://" + servers + "/admin?replicaSet=rs0"
	nums, err := strconv.Atoi(count)
	if err != nil {
		return err
	}
	if nums <= 0 {
		return fmt.Errorf("agrs count is not right!")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	collection := client.Database(db).Collection(coll)
	for num := 0; num < 100; num++ {
		ctx, _ := context.WithTimeout(context.Background(), 300*time.Second)
		uid, err := uuid.New()
		if err != nil {
			//return nil, err
			return err
		}
		res, err := collection.InsertOne(ctx, bson.M{"UUID": uid, "value": rand.Float32(), "timestamp": time.Now().String()})
		if err != nil {
			//return nil, err
			return err
		}
		fmt.Sprint(res.InsertedID)
	}

	return nil
}

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
