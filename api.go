package main

import (
	"fmt"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
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

func main() {
	db := Auth()

	// create a new table
	pk, _ := users.BuildPrimaryKey()
	usersTable := db.NewTable(users.TableName, pk)
	db.CreateTable(users)

	// load users data
	users := LoadUsersData()

	// put data int users table
	for i := range users {
		ok, err := usersTable.PutItem(strconv.FormatInt(int64(i), 10), "", users[i])
		if !ok {
			fmt.Println(err)
		}
	}
}
