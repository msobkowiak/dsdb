package main

import (
	"fmt"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"log"
)

var currentId int

var todos Todos
var tables Tables

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

func RepoCreateTable(t Table) Table {
	var tab = dynamodb.TableDescriptionT{
		TableName: t.Name,
		AttributeDefinitions: []dynamodb.AttributeDefinitionT{
			dynamodb.AttributeDefinitionT{"id", t.Id},
		},
		KeySchema: []dynamodb.KeySchemaT{
			dynamodb.KeySchemaT{"id", t.IdType},
		},
		ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
			ReadCapacityUnits:  t.ReadThrouput,
			WriteCapacityUnits: t.WriteThrouput,
		},
	}

	dynamodbRegion := aws.Region{DynamoDBEndpoint: "http://127.0.0.1:4567"}
	dynamodbAuth := aws.Auth{AccessKey: "key", SecretKey: "secret"}

	ddbs := dynamodb.Server{
		Auth:   dynamodbAuth,
		Region: dynamodbRegion,
	}

	// create a new table
	pk, _ := tab.BuildPrimaryKey()
	catalog := ddbs.NewTable(tab.TableName, pk)
	log.Println(catalog)
	ddbs.CreateTable(tab)

	tables = append(tables, t)
	return t
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
