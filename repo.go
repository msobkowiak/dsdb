package main

import (
	"github.com/goamz/goamz/dynamodb"
	//"log"
)

func RepoGetAllItems(tableName string) ([]map[string]string, error) {
	table := GetDynamoTable(tableName)
	items, err := table.Scan(nil)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemByHash(tableName, hash string) (map[string]string, error) {
	table := GetDynamoTable(tableName)
	item, err := table.GetItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return nil, err
	}

	return getData(item), nil
}

func RepoGetItemByHashRange(tableName, hashKey, rangeKey string) (map[string]string, error) {
	table := GetDynamoTable(tableName)
	item, err := table.GetItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		return nil, err
	}

	return getData(item), nil
}

func RepoGetItemByRange(tableName, hashValue string) ([]map[string]string, error) {
	table := GetDynamoTable(tableName)
	hashName := GetTableDescription(tableName).PrimaryKey.Hash

	atrrComaparations := buildQueryHash(tableName, hashName, hashValue)

	items, err := table.Query(atrrComaparations)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemByIndexHash(tableName, indexName, hashValue string) ([]map[string]string, error) {
	table := GetDynamoTable(tableName)
	schema := GetTableDescription(tableName)
	index, err := schema.GetIndexByName(indexName)
	if err != nil {
		return nil, err
	}

	atrrComaparations := buildQueryHash(tableName, index.Hash, hashValue)

	items, err := table.QueryOnIndex(atrrComaparations, indexName)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemsByRangeOp(tableName, hashValue, operator string, rangeValue []string) ([]map[string]string, error) {
	table := GetDynamoTable(tableName)
	schema := GetTableDescription(tableName)

	var atrrComaparations []dynamodb.AttributeComparison
	atrrComaparations = buildQueryRange(tableName, schema.PrimaryKey.Hash, hashValue, operator, schema.PrimaryKey.Range, rangeValue)

	items, err := table.Query(atrrComaparations)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemsByIndexRangeOp(tableName, indexName, hashValue, operator string, rangeValue []string) ([]map[string]string, error) {
	table := GetDynamoTable(tableName)
	schema := GetTableDescription(tableName)
	index, err := schema.GetIndexByName(indexName)
	if err != nil {
		return nil, err
	}

	var atrrComaparations []dynamodb.AttributeComparison
	atrrComaparations = buildQueryRange(tableName, index.Hash, hashValue, operator, index.Range, rangeValue)

	items, err := table.QueryOnIndex(atrrComaparations, indexName)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoDeleteItem(tableName, hash string) (bool, error) {
	schema := GetTableDescription(tableName)

	if schema.HasRange() {
		return deleteItems(tableName, hash, schema)
	}

	return deleteItem(tableName, hash)
}

func RepoDeleteItemWithRange(tableName, hashKey, rangeKey string) (bool, error) {
	table := GetDynamoTable(tableName)
	status, err := table.DeleteItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		return status, err
	}

	return status, nil
}

/*
[
	{"Name":"first_name", "Type":"S", "Value":"monika"},
	{"Name":"email", "Type":"S", "Value":"monika@example.pl"},
	{"Name":"last_name", "Type":"S", "Value":"Nowak"},
	{"Name":"counrty", "Type":"S", "Value":"Poland"}
]
*/
func RepoAddItem(tableName, hash string, item []Attribute) (bool, error) {
	table := GetDynamoTable(tableName)

	var attr = make([]dynamodb.Attribute, len(item))
	for i := range item {
		attr[i] = dynamodb.Attribute{
			Type:  item[i].Description.Type,
			Name:  item[i].Description.Name,
			Value: item[i].Value,
		}
	}
	return table.PutItem(hash, "", attr)
}

func RepoAddItemHashRange(tableName, hashKey, rangeKey string, item []Attribute) (bool, error) {
	table := GetDynamoTable(tableName)

	var attr = make([]dynamodb.Attribute, len(item))
	for i := range item {
		attr[i] = dynamodb.Attribute{
			Type:  item[i].Description.Type,
			Name:  item[i].Description.Name,
			Value: item[i].Value,
		}
	}
	return table.PutItem(hashKey, rangeKey, attr)
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

func buildQueryRange(tableName, hashName, hashValue, operator, rangeName string, rangeValue []string) []dynamodb.AttributeComparison {
	schema := GetTableDescription(tableName)

	var atrrs1 = make([]dynamodb.Attribute, 1)
	atrrs1[0] = dynamodb.Attribute{
		Value: hashValue,
		Name:  hashName,
		Type:  schema.GetTypeOfAttribute(hashName),
	}

	var atrrs2 []dynamodb.Attribute
	if len(rangeValue) == 1 {
		atrrs2 = make([]dynamodb.Attribute, 1)
		atrrs2[0] = dynamodb.Attribute{
			Value: rangeValue[0],
			Name:  rangeName,
		}
	} else if len(rangeValue) == 2 {
		atrrs2 = make([]dynamodb.Attribute, 2)
		atrrs2[0] = dynamodb.Attribute{
			Value: rangeValue[0],
			Name:  rangeName,
		}
		atrrs2[1] = dynamodb.Attribute{
			Value: rangeValue[1],
			Name:  rangeName,
		}
	}
	atrrs2[0] = dynamodb.Attribute{
		Value: rangeValue[0],
		Name:  rangeName,
		Type:  schema.GetTypeOfAttribute(rangeName),
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 2)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      hashName,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs1,
	}
	atrrComaparations[1] = dynamodb.AttributeComparison{
		AttributeName:      rangeName,
		ComparisonOperator: operator,
		AttributeValueList: atrrs2,
	}

	return atrrComaparations
}

func buildQueryHash(tableName, hashName, hashValue string) []dynamodb.AttributeComparison {
	schema := GetTableDescription(tableName)

	var atrrs = make([]dynamodb.Attribute, 1)
	atrrs[0] = dynamodb.Attribute{
		Value: hashValue,
		Name:  hashName,
		Type:  schema.GetTypeOfAttribute(hashName),
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 1)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      hashName,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs,
	}

	return atrrComaparations
}

func deleteItem(tableName, hash string) (bool, error) {
	table := GetDynamoTable(tableName)

	status, err := table.DeleteItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return status, err
	}

	return status, nil
}

func deleteItems(tableName, hash string, schema TableDescription) (bool, error) {
	table := GetDynamoTable(tableName)

	items, _ := RepoGetItemByRange(tableName, hash)
	for i := range items {
		status, err := table.DeleteItem(&dynamodb.Key{HashKey: items[i][schema.PrimaryKey.Hash], RangeKey: items[i][schema.PrimaryKey.Range]})
		if err != nil {
			return status, err
		}
	}

	return true, nil
}
