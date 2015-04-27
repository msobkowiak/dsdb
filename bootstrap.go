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

func deleteAllTables(db dynamodb.Server) {
	tables, err := db.ListTables()
	if err != nil {
		log.Println(err)
	} else {
		for i := range tables {
			deleteTable(db, tables[i])
		}
	}
}

func createUsersTable(db dynamodb.Server) {
	// create a new table
	tab := GetTableDescription(GetUsersSchema())
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

func createGameScoreTable(db dynamodb.Server) {
	// create a new table
	tab := GetTableDescription(GetGameScoreSchema())
	pk, _ := tab.BuildPrimaryKey()
	table := db.NewTable(tab.TableName, pk)
	_, err := db.CreateTable(tab)
	if err != nil {
		log.Println(err)
	}

	// load data
	data, rangeKey := LoadGameScoreData()
	// put data into table
	for i := range data {
		ok, err := table.PutItem(strconv.FormatInt(int64(i+1), 10), rangeKey[i], data[i])
		if !ok {
			log.Println(err)
		}
	}
}

func createAllTables(db dynamodb.Server) {
	createUsersTable(db)
	createGameScoreTable(db)
}

func Bootstrap() {
	db := Auth("http://127.0.0.1:4567", "key", "secret")

	// cleanup the database
	//deleteAllTables(db)
	createAllTables(db)
}
