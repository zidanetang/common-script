package main

import (
	cli "gopkg.in/urfave/cli.v2"
	"github.com/olekukonko/tablewriter"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"strconv"

	"math/rand"
	"os"
	"context"
	"time"
	"fmt"
)



func PrintError(err error) {
	color.Set(color.FgRed)
	defer color.Unset()
	fmt.Printf("ERROR: %s\n", err)
}

func SetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:"servers",
			Aliases: []string{"s"},
			Usage: "Specify MongoDB server Cluster list, example: ['127.0.0.1:27001', '127.0.0.1:27002','127.0.0.1:27003']",
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

func PrintWithTable(title string, rows [][]string, header []string) {
	fmt.Printf(title)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	var headerColors []tablewriter.Colors
	for i := 0; i < len(header); i++ {
		headerColors = append(headerColors, tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor})
	}

	table.SetHeaderColor(headerColors...)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment

	for _, row := range rows {
		table.Append(row)
	}

	table.Render() // Send output
}

func Clinet(servers []string) (*mongo.Client, error) {
	var client *mongo.Client
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	for _, s := range servers{
		c, err := mongo.Connect(ctx,options.Client().ApplyURI("mongodb://"+ s))
		if err != nil {
			return nil, nil
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			return nil, nil
		}
		client = c
	}
	return client, nil
}

func insertDocuments(c *mongo.Client, count string, db string, table string) ([]string, error) {
	nums, err := strconv.Atoi(count)
	if err != nil {
		return  nil, err
	}
	if nums <= 0{
		return nil, fmt.Errorf("agrs count is not right!")
	}
	collection := c.Database(db).Collection(table)
	result := make([]string, nums)
	for num := 0; num < nums ; num++  {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		res, err := collection.InsertOne(ctx, bson.M{"UUID": uuid.New(),"value": rand.Float32(),"timestamp": time.Now().String()})
		if err != nil {
			return nil, err
		}
		//id := res.InsertedID
		record := fmt.Sprint( res.InsertedID)
		result = append(result, record)
	}
	return result, nil
}

func run(c *cli.Context) error {
	servers := c.StringSlice("servers")
	count := c.String("counts")
	db := c.String("DB")
	coll := c.String("Collection")
	client, err := Clinet(servers)
	if err != nil {
		return err
	}
	result, err := insertDocuments(client, count, db, coll)
	if err != nil || len(result) == 0 {
		return err
	}
	title := "Result:\n"
	 cont := [][]string{result}
	header := []string{"InsertedID"}

	PrintWithTable(title, cont, header)
	return nil
}

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
		Name:     "are",
		Usage:    "Auto Remediate AWS EC2 Events",
		Flags:    SetFlags(),
		Compiled: time.Now(),
		Version:  VERSION,
		Action: func(c *cli.Context) error {
			return run(c)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		PrintError(err)
	}
}
