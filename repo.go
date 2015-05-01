package main

import (
	"github.com/goamz/goamz/dynamodb"
	"log"
)

func RepoGetAllItems(tableName string) ([]map[string]string, error) {
	table := GetTable(tableName)
	items, err := table.Scan(nil)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemByHash(tableName, hash string) map[string]string {
	table := GetTable(tableName)
	item, err := table.GetItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		var ex = make(map[string]string)
		ex["exception"] = "The hash " + hash + " not found in table " + tableName
		log.Println(err)

		return ex
	}

	return getData(item)
}

func RepoDeleteItemByHash(tableName, hash string) (bool, error) {
	table := GetTable(tableName)
	status, err := table.DeleteItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return status, err
	}

	return status, nil
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

func RepoGetItemByRange(tableName, rangeKey, operator, value string) []map[string]*dynamodb.Attribute {
	//func RepoGetItemByHashRangeOp(tableName, hashKey, rangeKey, value string) map[string]string {
	table := GetTable(tableName)
	//item, err := table.GetItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	var attrToGet = dynamodb.Attribute{
		Value: value,
		Name:  rangeKey,
		Type:  "S",
	}

	var atrrs = make([]dynamodb.Attribute, 1)
	atrrs[0] = attrToGet

	var atrrComaration = dynamodb.AttributeComparison{
		AttributeName:      rangeKey,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs,
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 1)
	atrrComaparations[0] = atrrComaration
	item, err := table.Query(atrrComaparations)
	if err != nil {
		//var ex = make(map[string]string)
		//ex["exception"] = "The range " + rangeKey + " not found in table " + tableName
		log.Println(err)
		//return ex
	}
	log.Println(item)
	/*var data = make(map[string]string)
	for key := range item {
		data[key] = item[key].Value
	}*/

	return item
}

func getData(item map[string]*dynamodb.Attribute) map[string]string {
	var data = make(map[string]string)
	for key := range item {
		data[key] = item[key].Value
	}

	return data
}

func getDataAsArray(items []map[string]*dynamodb.Attribute) []map[string]string {
	itemCount := len(items)
	var data = make([]map[string]string, itemCount)
	for k := 0; k < itemCount; k++ {
		data[k] = getData(items[k])
	}

	return data
}
