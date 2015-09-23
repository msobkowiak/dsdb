package main

import (
	"log"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
)

type DynamoClient struct {
}

func (c *DynamoClient) Auth(region, accessKey, secretKey string) dynamodb.Server {
	dynamodbRegion := aws.Region{DynamoDBEndpoint: region}
	dynamodbAuth := aws.Auth{AccessKey: accessKey, SecretKey: secretKey}

	return dynamodb.Server{
		Auth:   dynamodbAuth,
		Region: dynamodbRegion,
	}
}

func (c *DynamoClient) DeleteAll(db dynamodb.Server) {
	var t DynamoTable
	tables, err := db.ListTables()
	if err != nil {
		log.Println(err)
	} else {
		for i := range tables {
			t.Delete(db, tables[i])
		}
	}
}
