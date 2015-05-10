package main

import (
	//"fmt"
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

func ConvertToDynamo(t TableDescription) dynamodb.TableDescriptionT {

	var table = dynamodb.TableDescriptionT{
		TableName: t.Name,
	}

	addAttrubuteDefinition(&table, t.Attributes)
	addPrimaryKey(&table, t.PrimaryKey, t.Attributes)
	addSecondaryIndexes(&table, t.SecondaryIndexes, t.Attributes)
	addThroughput(&table)

	return table
}

func addHash(hashName string, atrrs []AttributeDefinition, keySchema []dynamodb.KeySchemaT) {
	if existsInAttributes(atrrs, hashName) {
		keySchema[0] = dynamodb.KeySchemaT{hashName, "HASH"}
	}
}

func addRange(rangeName string, atrrs []AttributeDefinition, keySchema []dynamodb.KeySchemaT) {
	if existsInAttributes(atrrs, rangeName) {
		keySchema[1] = dynamodb.KeySchemaT{rangeName, "RANGE"}
	}
}

func existsInAttributes(attrs []AttributeDefinition, keyName string) bool {
	for i := range attrs {
		if attrs[i].Name == keyName {
			return true
		}
	}

	return false
}

func addAttrubuteDefinition(table *dynamodb.TableDescriptionT, attrs []AttributeDefinition) {
	table.AttributeDefinitions = make([]dynamodb.AttributeDefinitionT, len(attrs))
	for i := range attrs {
		table.AttributeDefinitions[i] = dynamodb.AttributeDefinitionT{attrs[i].Name, attrs[i].Type}
	}
}

func addPrimaryKey(table *dynamodb.TableDescriptionT, key PrimaryKeyDefinition, attrs []AttributeDefinition) {
	if key.Type == "HASH" {
		table.KeySchema = make([]dynamodb.KeySchemaT, 1)
		addHash(key.Hash, attrs, table.KeySchema)
	} else if key.Type == "RANGE" {
		table.KeySchema = make([]dynamodb.KeySchemaT, 2)
		addHash(key.Hash, attrs, table.KeySchema)
		addRange(key.Range, attrs, table.KeySchema)
	}
}

func addSecondaryIndexes(table *dynamodb.TableDescriptionT, indexes []SecondaryIndexDefinition, attrs []AttributeDefinition) {
	table.GlobalSecondaryIndexes = make([]dynamodb.GlobalSecondaryIndexT, len(indexes))

	for i := range indexes {

		table.GlobalSecondaryIndexes[i] = dynamodb.GlobalSecondaryIndexT{
			IndexName: indexes[i].Name,
		}

		if indexes[i].Type == "HASH" {
			table.GlobalSecondaryIndexes[i].KeySchema = make([]dynamodb.KeySchemaT, 1)
			addHash(indexes[i].Hash, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
		} else if indexes[i].Type == "RANGE" {
			table.GlobalSecondaryIndexes[i].KeySchema = make([]dynamodb.KeySchemaT, 2)
			addHash(indexes[i].Hash, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
			addRange(indexes[i].Range, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
		}
		table.GlobalSecondaryIndexes[i].Projection = dynamodb.ProjectionT{"ALL"}
		table.GlobalSecondaryIndexes[i].ProvisionedThroughput = dynamodb.ProvisionedThroughputT{
			ReadCapacityUnits:  1,
			WriteCapacityUnits: 1,
		}
	}
}

func addThroughput(table *dynamodb.TableDescriptionT) {
	table.ProvisionedThroughput = dynamodb.ProvisionedThroughputT{
		ReadCapacityUnits:  10,
		WriteCapacityUnits: 10,
	}
}

func GetTable(tableName string) dynamodb.Table {
	db := Auth("http://127.0.0.1:4567", "key", "secret")

	tableDescription := ConvertToDynamo(GetSchema(tableName))
	pk, _ := tableDescription.BuildPrimaryKey()

	return *db.NewTable(tableDescription.TableName, pk)
}
