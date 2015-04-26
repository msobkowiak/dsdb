package main

import (
	"github.com/goamz/goamz/dynamodb"
	"log"
	"strconv"
)

func deleteTable(db dynamodb.Server, table string) {
	tabDescription, err := db.DescribeTable(table)
	if err != nil {
		log.Println(err)
	} else {
		_, err := db.DeleteTable(*tabDescription)
		if err != nil {
			log.Println(err)
		}
	}
}

func deleteAllTables(db dynamodb.Server, tables []string) {
	for i := range tables {
		deleteTable(db, tables[i])
	}
}

func createTable(db dynamodb.Server, tab dynamodb.TableDescriptionT) {
	// create a new table
	pk, _ := tab.BuildPrimaryKey()
	table := db.NewTable(tab.TableName, pk)
	_, err := db.CreateTable(tab)
	if err != nil {
		log.Println(err)
	}

	// load data
	data := LoadUsersData()
	// put data into table
	for i := range data {
		ok, err := table.PutItem(strconv.FormatInt(int64(i+1), 10), "", data[i])
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
	db := Auth("http://127.0.0.1:4567", "key", "secret")

	// cleanup the database
	/*var tables = make([]string, 1)
	tables[0] = "users"
	deleteAllTables(db, tables)*/

	// create tables
	var tablesDescription = make([]dynamodb.TableDescriptionT, 1)
	tablesDescription[0] = GetTableDescription(GetUsersSchema())
	createAllTables(db, tablesDescription)
}
