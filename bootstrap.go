package main

import (
	"github.com/goamz/goamz/dynamodb"
	"log"
	"strconv"
	"time"
)

const TIMEOUT = 1 * time.Minute

func deleteTable(db dynamodb.Server, tableName string) {
	tabDescription, err := db.DescribeTable(tableName)
	if err != nil {
		log.Println(err)
	} else {
		_, err := db.DeleteTable(*tabDescription)
		if err != nil {
			log.Println(err)
		}
	}

	pk, _ := tabDescription.BuildPrimaryKey()
	table := db.NewTable(tableName, pk)
	WaitUntilTableDeleted(db, table, tableName)
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

func createTable(db dynamodb.Server, tableName string) {
	// create a new table
	tab := ConvertToDynamo(GetTableDescription(tableName))
	pk, _ := tab.BuildPrimaryKey()
	table := db.NewTable(tab.TableName, pk)
	_, err := db.CreateTable(tab)
	if err != nil {
		log.Println(err)
	}

	// load data
	data, hashKeys := GetTableData(tableName)

	// put data into table
	WaitUntilStatus(table, "ACTIVE")
	if hashKeys != nil {
		for i := range data {
			ok, err := table.PutItem(hashKeys[i], strconv.FormatInt(int64(i+1), 10), data[i])
			if !ok {
				log.Println(err)
			}
		}
	} else {
		for i := range data {
			ok, err := table.PutItem(strconv.FormatInt(int64(i+1), 10), "", data[i])
			if !ok {
				log.Println(err)
			}
		}
	}

}

func WaitUntilTableDeleted(db dynamodb.Server, t *dynamodb.Table, tableName string) {
	done := make(chan bool)
	timeout := time.After(TIMEOUT)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				tables, err := db.ListTables()
				if err != nil {
					log.Fatal(err)
				}
				if findTableByName(tables, tableName) {
					time.Sleep(5 * time.Second)
				} else {
					done <- true
					return
				}
			}
		}
	}()
	select {
	case <-done:
		break
	case <-timeout:
		log.Println("Expect the table to be deleted but timed out")
		close(done)
	}
}

func WaitUntilStatus(t *dynamodb.Table, status string) {
	// We should wait until the table is in specified status because a real DynamoDB has some delay for ready
	done := make(chan bool)
	timeout := time.After(TIMEOUT)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				desc, err := t.DescribeTable()
				if err != nil {
					log.Fatal(err)
				}
				if desc.TableStatus == status {
					done <- true
					return
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()
	select {
	case <-done:
		break
	case <-timeout:
		log.Printf("Expect a status to be %s, but timed out\n", status)
		close(done)
	}
}

func findTableByName(tables []string, name string) bool {
	for _, t := range tables {
		if t == name {
			return true
		}
	}
	return false
}

func createAllTables(db dynamodb.Server) {
	createTable(db, "users")
	createTable(db, "game_scores")
}

func Bootstrap() {
	db := Auth("http://127.0.0.1:4567", "key", "secret")

	// cleanup the database
	deleteAllTables(db)
	createAllTables(db)
}
