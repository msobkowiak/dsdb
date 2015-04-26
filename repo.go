package main

import (
	"fmt"
	"github.com/goamz/goamz/dynamodb"
)

var currentId int

var todos Todos

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

func RepoGetItemByHash(tableName string, hash string) map[string]string {
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

func RepoDestroyTodo(id int) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
