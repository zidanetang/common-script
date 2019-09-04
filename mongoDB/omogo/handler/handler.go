package handler

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	cli "gopkg.in/urfave/cli.v2"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"
)

func PrintError(err error) {
	color.Set(color.FgRed)
	defer color.Unset()
	fmt.Printf("ERROR: %s\n", err)
}

func SetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "servers",
			Aliases: []string{"s"},
			Usage: "Specify MongoDB server Cluster, example, e.g.\n" +
				"\t\t\t -s 127.0.0.1:27017\n",
			//"\t\t\t -s \"127.0.0.1:27017 127.0.0.1:27016 127.0.0.1:27015\"",
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

func Duplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}
func Clinet(servers string) (*mongo.Client, error) {
	//func Clinet(servers []string) (*mongo.Client, error) {
	var client *mongo.Client
	ctx, _ := context.WithTimeout(context.Background(), 300*time.Second)
	/*
		for _, s := range servers {
			c, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+s))
			if err != nil {
				return nil, nil
			}
			err = client.Ping(ctx, readpref.Primary())
			if err != nil {
				return nil, nil
			}
			client = c
		}
	*/
	server := "mongodb://" + servers
	fmt.Println(server)
	c, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+server))
	if err != nil {
		return nil, nil
	}
	client = c
	return client, nil
}

func insertDocuments(c *mongo.Client, count string, db string, table string) ([]string, error) {
	nums, err := strconv.Atoi(count)
	if err != nil {
		return nil, err
	}
	if nums <= 0 {
		return nil, fmt.Errorf("agrs count is not right!")
	}
	collection := c.Database(db).Collection(table)
	result := make([]string, nums)
	for num := 0; num < nums; num++ {
		ctx, _ := context.WithTimeout(context.Background(), 300*time.Second)
		uid, err := uuid.New()
		if err != nil {
			return nil, err
		}
		res, err := collection.InsertOne(ctx, bson.M{"UUID": uid, "value": rand.Float32(), "timestamp": time.Now().String()})
		if err != nil {
			return nil, err
		}
		//id := res.InsertedID
		record := fmt.Sprint(res.InsertedID)
		result = append(result, record)
	}
	return result, nil
}

func Run(c *cli.Context) error {
	//var serviceTodoList []string
	servers := c.String("servers")
	/*
		if len(servers) != 0 {
			reg := regexp.MustCompile(`\s+`)
			serversList := reg.Split(strings.TrimSpace(servers), -1)
			sort.Strings(serversList)
			distinctServersList := Duplicate(serversList)
			for _, service := range distinctServersList {
				serviceTodoList = append(serviceTodoList, service.(string))
			}
		}
	*/
	count := c.String("counts")
	db := c.String("DB")
	coll := c.String("Collection")
	client, err := Clinet(servers)
	//client, err := Clinet(serviceTodoList)
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
