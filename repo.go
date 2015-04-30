package main

import (
	"fmt"
	"github.com/goamz/goamz/dynamodb"
	//"log"
)

var currentId int

var todos Todos

type fetchedItems struct {
	item map[string]string
}

// Give us some seed data
func init() {
	RepoCreateTodo(Todo{Name: "Write presentation"})
	RepoCreateTodo(Todo{Name: "Host meetup"})
}

func RepoFindTodo(id int) Todo {
	for _, t := range todos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return Todo{}
}

func RepoCreateTodo(t Todo) Todo {
	currentId += 1
	t.Id = currentId
	todos = append(todos, t)
	return t
}

func RepoGetAllItems(tableName string) []map[string]string {
	table := GetTable(tableName)
	items, err := table.Scan(nil)
	if err != nil {
		var ex = make([]map[string]string, 1)
		var exception = make(map[string]string)
		exception["exception"] = "The table " + tableName + " not found"
		ex[0] = exception
		return ex
	}

	itemCount := len(items)
	var data = make([]map[string]string, itemCount)
	for k := 0; k < itemCount; k++ {
		item := make(map[string]string)
		for key := range items[k] {
			item[key] = items[k][key].Value
		}
		data[k] = item
	}

	return data
}

func RepoGetItemByHash(tableName, hash string) map[string]string {
	table := GetTable(tableName)
	item, err := table.GetItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		var ex = make(map[string]string)
		ex["exception"] = "The hash " + hash + " not found in table " + tableName

		return ex
	}

	var data = make(map[string]string)
	for key := range item {
		data[key] = item[key].Value
	}

	return data
}

func RepoGetItemByHashRange(tableName, hashKey, rangeKey string) map[string]string {
	table := GetTable(tableName)
	item, err := table.GetItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		var ex = make(map[string]string)
		ex["exception"] = "The hash " + hashKey + " and range " + rangeKey + " not found in table " + tableName

		return ex
	}

	var data = make(map[string]string)
	for key := range item {
		data[key] = item[key].Value
	}

	return data
}

func RepoGetItemByHashRangeOp(tableName, hashKey, rangeKey, operator, value string) map[string]string {
	//func RepoGetItemByHashRangeOp(tableName, hashKey, rangeKey, value string) map[string]string {
	table := GetTable(tableName)
	//log.Println(operator)
	item, err := table.GetItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		var ex = make(map[string]string)
		ex["exception"] = "AAAA The hash " + hashKey + " and range " + rangeKey + " not found in table " + tableName

		return ex
	}

	var data = make(map[string]string)
	for key := range item {
		data[key] = item[key].Value
	}

	return data
}

func RepoDestroyTodo(id int) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
