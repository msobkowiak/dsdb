package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"log"
	"strconv"
)

func Auth() dynamodb.Server {
	dynamodbRegion := aws.Region{DynamoDBEndpoint: "http://127.0.0.1:4567"}
	dynamodbAuth := aws.Auth{AccessKey: "key", SecretKey: "secret"}

	return dynamodb.Server{
		Auth:   dynamodbAuth,
		Region: dynamodbRegion,
	}
}

func deleteTable(db dynamodb.Server, table string) {
	tabDescription, err := db.DescribeTable(table)
	if err != nil {
		delete, _ := db.DeleteTable(*tabDescription)
		log.Println(delete)
	}
}

func deleteAllTables(db dynamodb.Server, tables []string) {
	for i := range tables {
		deleteTable(db, tables[i])
	}
}

func createTable(db dynamodb.Server, tab dynamodb.TableDescriptionT) {
	// create a new table
	pk, _ := users.BuildPrimaryKey()
	table := db.NewTable(tab.TableName, pk)
	db.CreateTable(tab)

	// load data
	data := LoadData()

	// put data into table
	for i := range data {
		ok, err := table.PutItem(strconv.FormatInt(int64(i), 10), "", data[i])
		if !ok {
			log.Println(err)
		}
	}
}

func createAllTables(db dynamodb.Server, tables []dynamodb.TableDescriptionT) {
	for i := range tables {
		createTable(db, tables[i])
	}
}

func Bootstrap() {
	db := Auth()

	// cleanup the database
	var tables = make([]string, 1)
	tables[0] = "Users"
	deleteAllTables(db, tables)

	// create tables
	var tablesDescription = make([]dynamodb.TableDescriptionT, 1)
	tablesDescription[0] = users
	createAllTables(db, tablesDescription)
}
