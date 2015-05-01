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
func GetTableDescription(t Table) (dynamodb.TableDescriptionT, bool) {
	//if t.HasRange() {

	if t.HasRange() {
		return getTableDescriptionHashRange(t), false
	}
	return getTableDescriptionHash(t), true
}

func getTableDescriptionHash(t Table) dynamodb.TableDescriptionT {
	return dynamodb.TableDescriptionT{
		TableName: t.Name,
		AttributeDefinitions: []dynamodb.AttributeDefinitionT{
			dynamodb.AttributeDefinitionT{t.HashKey.Name, t.HashKey.AttributeType},
		},
		KeySchema: []dynamodb.KeySchemaT{
			dynamodb.KeySchemaT{t.HashKey.Name, t.HashKey.KeyType},
		},
		ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
			ReadCapacityUnits:  t.ReadCapacityUnits,
			WriteCapacityUnits: t.WriteCapacityUnits,
		},
	}
}

func getTableDescriptionHashRange(t Table) dynamodb.TableDescriptionT {
	return dynamodb.TableDescriptionT{
		TableName: t.Name,
		AttributeDefinitions: []dynamodb.AttributeDefinitionT{
			dynamodb.AttributeDefinitionT{t.HashKey.Name, t.HashKey.AttributeType},
			dynamodb.AttributeDefinitionT{t.RangeKey.Name, t.RangeKey.AttributeType},
		},
		KeySchema: []dynamodb.KeySchemaT{
			dynamodb.KeySchemaT{t.HashKey.Name, t.HashKey.KeyType},
			dynamodb.KeySchemaT{t.RangeKey.Name, t.RangeKey.KeyType},
		},
		ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
			ReadCapacityUnits:  t.ReadCapacityUnits,
			WriteCapacityUnits: t.WriteCapacityUnits,
		},
	}
}

func GetTable(tableName string) (dynamodb.Table, bool) {
	db := Auth("http://127.0.0.1:4567", "key", "secret")

	tableDescription, hasRange := GetTableDescription(GetSchema(tableName))
	pk, _ := tableDescription.BuildPrimaryKey()

	return *db.NewTable(tableDescription.TableName, pk), hasRange
}
