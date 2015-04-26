package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
)

// "http://127.0.0.1:4567"
// "key"
// "secret"

func Auth(region, accessKey, secretKey string) dynamodb.Server {
	dynamodbRegion := aws.Region{DynamoDBEndpoint: region}
	dynamodbAuth := aws.Auth{AccessKey: accessKey, SecretKey: secretKey}

	return dynamodb.Server{
		Auth:   dynamodbAuth,
		Region: dynamodbRegion,
	}
}

func GetTableDescription(t Table) dynamodb.TableDescriptionT {
	return dynamodb.TableDescriptionT{
		TableName: t.Name,
		AttributeDefinitions: []dynamodb.AttributeDefinitionT{
			dynamodb.AttributeDefinitionT{t.PrimaryKey.Name, t.PrimaryKey.AttributeType},
		},
		KeySchema: []dynamodb.KeySchemaT{
			dynamodb.KeySchemaT{t.PrimaryKey.Name, t.PrimaryKey.KeyType},
		},
		ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
			ReadCapacityUnits:  t.ReadCapacityUnits,
			WriteCapacityUnits: t.WriteCapacityUnits,
		},
	}
}

func GetTable(tableName string) dynamodb.Table {
	db := Auth("http://127.0.0.1:4567", "key", "secret")

	tableDescription := GetTableDescription(GetUsersSchema())
	pk, _ := tableDescription.BuildPrimaryKey()

	return *db.NewTable(tableDescription.TableName, pk)
}
