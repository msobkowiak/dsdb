package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"log"
	"strconv"
	"time"
)

const TIMEOUT = 1 * time.Minute

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

	addAttributeDefinition(&table, t)
	addPrimaryKey(&table, t.PrimaryKey, t.Attributes)
	addSecondaryIndexes(&table, t.SecondaryIndexes, t.Attributes)
	addThroughput(&table)

	return table
}

func GetDynamoTable(tableName string) (dynamodb.Table, error) {
	tabDesc, err := GetTableDescription(tableName, schema.Tables)
	if err != nil {
		return dynamodb.Table{}, err
	}

	auth := schema.Authentication.Dynamo
	db := Auth(auth.Region, auth.AccessKey, auth.SecretKey)

	dynamoTab := ConvertToDynamo(tabDesc)
	pk, _ := dynamoTab.BuildPrimaryKey()

	return *db.NewTable(dynamoTab.TableName, pk), nil
}

func DeleteAllTables(db dynamodb.Server) {
	tables, err := db.ListTables()
	if err != nil {
		log.Println(err)
	} else {
		for i := range tables {
			deleteTable(db, tables[i])
		}
	}
}

func CreateTable(t TableDescription) dynamodb.Table {
	// get dynamoDB Server
	dynamAuth := schema.Authentication.Dynamo
	db := Auth(dynamAuth.Region, dynamAuth.AccessKey, dynamAuth.SecretKey)

	// create a new table
	tab := ConvertToDynamo(t)
	pk, _ := tab.BuildPrimaryKey()
	table := db.NewTable(tab.TableName, pk)
	_, err := db.CreateTable(tab)
	if err != nil {
		log.Println(err)
	}
	waitUntilStatus(table, "ACTIVE")

	return *table
}

func AddItems(t dynamodb.Table, data [][]dynamodb.Attribute, hashKeys []string) {
	if hashKeys != nil {
		for i := range data {
			ok, err := t.PutItem(hashKeys[i], strconv.FormatInt(int64(i+1), 10), data[i])
			if !ok {
				log.Println(err)
			}
		}
	} else {
		for i := range data {
			ok, err := t.PutItem(strconv.FormatInt(int64(i+1), 10), "", data[i])
			if !ok {
				log.Println(err)
			}
		}
	}
}

func waitUntilTableDeleted(db dynamodb.Server, t *dynamodb.Table, tableName string) {
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

func waitUntilStatus(t *dynamodb.Table, status string) {
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
	waitUntilTableDeleted(db, table, tableName)
}

func findTableByName(tables []string, name string) bool {
	for _, t := range tables {
		if t == name {
			return true
		}
	}
	return false
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

func addAttributeDefinition(table *dynamodb.TableDescriptionT, schemaTable TableDescription) {
	attrs := ExcludeNonKeyAttributes(schemaTable)
	table.AttributeDefinitions = make([]dynamodb.AttributeDefinitionT, len(attrs))
	for i := range attrs {
		table.AttributeDefinitions[i] = dynamodb.AttributeDefinitionT{attrs[i].Name, attrs[i].Type}
	}
}

func addPrimaryKey(table *dynamodb.TableDescriptionT, key KeyDefinition, attrs []AttributeDefinition) {
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

		if indexes[i].Key.Type == "HASH" {
			table.GlobalSecondaryIndexes[i].KeySchema = make([]dynamodb.KeySchemaT, 1)
			addHash(indexes[i].Key.Hash, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
		} else if indexes[i].Key.Type == "RANGE" {
			table.GlobalSecondaryIndexes[i].KeySchema = make([]dynamodb.KeySchemaT, 2)
			addHash(indexes[i].Key.Hash, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
			addRange(indexes[i].Key.Range, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
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
