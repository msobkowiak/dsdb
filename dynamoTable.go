package main

import (
	"github.com/goamz/goamz/dynamodb"
	"log"
	"strconv"
	"time"
)

const TIMEOUT = 1 * time.Minute

type DynamoTable struct {
}

func (d DynamoTable) Create(t TableDescription) dynamodb.Table {

	var client DynamoClient
	// get dynamoDB Server
	dynamAuth := schema.Authentication.Dynamo
	db := client.Auth(dynamAuth.Region, dynamAuth.AccessKey, dynamAuth.SecretKey)

	// create a new table
	tab := d.Map(t)
	pk, _ := tab.BuildPrimaryKey()
	table := db.NewTable(tab.TableName, pk)
	_, err := db.CreateTable(tab)
	if err != nil {
		log.Println(err)
	}
	d.waitUntilStatus(table, "ACTIVE")

	return *table
}

func (d DynamoTable) AddItems(t dynamodb.Table, data [][]dynamodb.Attribute, hashKeys []string) {
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

func (d DynamoTable) GetByName(tableName string) (dynamodb.Table, error) {
	var client DynamoClient
	tabDesc, err := schema.GetTableDescription(tableName)
	if err != nil {
		return dynamodb.Table{}, err
	}

	auth := schema.Authentication.Dynamo
	db := client.Auth(auth.Region, auth.AccessKey, auth.SecretKey)

	dynamoTab := d.Map(tabDesc)
	pk, _ := dynamoTab.BuildPrimaryKey()

	return *db.NewTable(dynamoTab.TableName, pk), nil
}

func (d DynamoTable) Delete(db dynamodb.Server, tableName string) {
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
	d.waitUntilTableDeleted(db, table, tableName)
}

func (d DynamoTable) waitUntilTableDeleted(db dynamodb.Server, t *dynamodb.Table, tableName string) {
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
				if d.findTableByName(tables, tableName) {
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

func (d DynamoTable) waitUntilStatus(t *dynamodb.Table, status string) {
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

func (d DynamoTable) findTableByName(tables []string, name string) bool {
	for _, t := range tables {
		if t == name {
			return true
		}
	}
	return false
}

func (d DynamoTable) addHash(hashName string, atrrs []AttributeDefinition, keySchema []dynamodb.KeySchemaT) {
	if existsInAttributes(atrrs, hashName) {
		keySchema[0] = dynamodb.KeySchemaT{hashName, "HASH"}
	}
}

func (d DynamoTable) addRange(rangeName string, atrrs []AttributeDefinition, keySchema []dynamodb.KeySchemaT) {
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

func (d DynamoTable) addAttributeDefinition(table *dynamodb.TableDescriptionT, schemaTable TableDescription) {
	attrs := ExcludeNonKeyAttributes(schemaTable)
	table.AttributeDefinitions = make([]dynamodb.AttributeDefinitionT, len(attrs))
	for i := range attrs {
		table.AttributeDefinitions[i] = dynamodb.AttributeDefinitionT{attrs[i].Name, attrs[i].Type}
	}
}

func (d DynamoTable) addPrimaryKey(table *dynamodb.TableDescriptionT, key KeyDefinition, attrs []AttributeDefinition) {
	if key.Type == "HASH" {
		table.KeySchema = make([]dynamodb.KeySchemaT, 1)
		d.addHash(key.Hash, attrs, table.KeySchema)
	} else if key.Type == "RANGE" {
		table.KeySchema = make([]dynamodb.KeySchemaT, 2)
		d.addHash(key.Hash, attrs, table.KeySchema)
		d.addRange(key.Range, attrs, table.KeySchema)
	}
}

func (d DynamoTable) addSecondaryIndexes(table *dynamodb.TableDescriptionT, indexes []SecondaryIndexDefinition, attrs []AttributeDefinition) {
	table.GlobalSecondaryIndexes = make([]dynamodb.GlobalSecondaryIndexT, len(indexes))

	for i := range indexes {

		table.GlobalSecondaryIndexes[i] = dynamodb.GlobalSecondaryIndexT{
			IndexName: indexes[i].Name,
		}

		if indexes[i].Key.Type == "HASH" {
			table.GlobalSecondaryIndexes[i].KeySchema = make([]dynamodb.KeySchemaT, 1)
			d.addHash(indexes[i].Key.Hash, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
		} else if indexes[i].Key.Type == "RANGE" {
			table.GlobalSecondaryIndexes[i].KeySchema = make([]dynamodb.KeySchemaT, 2)
			d.addHash(indexes[i].Key.Hash, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
			d.addRange(indexes[i].Key.Range, attrs, table.GlobalSecondaryIndexes[i].KeySchema)
		}
		table.GlobalSecondaryIndexes[i].Projection = dynamodb.ProjectionT{"ALL"}
		table.GlobalSecondaryIndexes[i].ProvisionedThroughput = dynamodb.ProvisionedThroughputT{
			ReadCapacityUnits:  1,
			WriteCapacityUnits: 1,
		}
	}
}

func (d DynamoTable) addThroughput(table *dynamodb.TableDescriptionT) {
	table.ProvisionedThroughput = dynamodb.ProvisionedThroughputT{
		ReadCapacityUnits:  10,
		WriteCapacityUnits: 10,
	}
}

func (d DynamoTable) Map(t TableDescription) dynamodb.TableDescriptionT {

	var table = dynamodb.TableDescriptionT{
		TableName: t.Name,
	}

	d.addAttributeDefinition(&table, t)
	d.addPrimaryKey(&table, t.PrimaryKey, t.Attributes)
	d.addSecondaryIndexes(&table, t.SecondaryIndexes, t.Attributes)
	d.addThroughput(&table)

	return table
}
