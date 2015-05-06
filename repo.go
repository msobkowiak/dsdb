package main

import (
	"github.com/goamz/goamz/dynamodb"
	//"log"
)

func RepoGetAllItems(tableName string) ([]map[string]string, error) {
	table, _ := GetTable(tableName)
	items, err := table.Scan(nil)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemByHash(tableName, hash string) (map[string]string, error) {
	table, _ := GetTable(tableName)

	item, err := table.GetItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return nil, err
	}

	return getData(item), nil
}

func RepoGetItemByHashRange(tableName, hashKey, rangeKey string) (map[string]string, error) {
	table, _ := GetTable(tableName)
	item, err := table.GetItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		return nil, err
	}

	return getData(item), nil
}

func RepoGetItemByRange(tableName, hash string) ([]map[string]string, error) {
	table, _ := GetTable(tableName)
	atrrComaparations := buildQueryHash(tableName, hash)

	items, err := table.Query(atrrComaparations)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemByIndexHash(tableName, hash string) ([]map[string]string, error) {
	table, _ := GetTable(tableName)
	schema := GetSchema(tableName)

	atrrComaparations := buildQueryIndexHash(tableName, hash)

	items, err := table.QueryOnIndex(atrrComaparations, schema.GlobalSecondaryIndex.Name)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemByIndexRange(tableName, hash string) ([]map[string]string, error) {
	table, _ := GetTable(tableName)
	schema := GetSchema(tableName)
	atrrComaparations := buildQueryIndexHashRange(tableName, hash)

	items, err := table.QueryOnIndex(atrrComaparations, schema.GlobalSecondaryIndex.Name)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemsByRangeOp(tableName, hash, operator string, rangeValue []string) ([]map[string]string, error) {
	table, _ := GetTable(tableName)

	var atrrComaparations []dynamodb.AttributeComparison
	if len(rangeValue) == 1 {
		atrrComaparations = buildQueryRange(tableName, hash, operator, rangeValue[0])
	} else {
		atrrComaparations = buildQueryRangeBetween(tableName, hash, operator, rangeValue)
	}

	items, err := table.Query(atrrComaparations)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoGetItemsByIndexRangeOp(tableName, hash, operator string, rangeValue []string) ([]map[string]string, error) {
	table, _ := GetTable(tableName)
	schema := GetSchema(tableName)

	var atrrComaparations []dynamodb.AttributeComparison
	if len(rangeValue) == 1 {
		atrrComaparations = buildQueryIndexRange(tableName, hash, operator, rangeValue[0])
	} else {
		atrrComaparations = buildQueryIndexRangeBetween(tableName, hash, operator, rangeValue)
	}

	items, err := table.QueryOnIndex(atrrComaparations, schema.GlobalSecondaryIndex.Name)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func RepoDeleteItem(tableName, hash string) (bool, error) {
	schema := GetSchema(tableName)

	if schema.HasRange() {
		return deleteItems(tableName, hash, schema)
	}

	return deleteItem(tableName, hash)
}

func RepoDeleteItemWithRange(tableName, hashKey, rangeKey string) (bool, error) {
	table, _ := GetTable(tableName)
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
	table, _ := GetTable(tableName)

	var attr = make([]dynamodb.Attribute, len(item))
	for i := range item {
		attr[i] = dynamodb.Attribute{
			Type:  item[i].Type,
			Name:  item[i].Name,
			Value: item[i].Value,
		}
	}
	return table.PutItem(hash, "", attr)
}

func RepoAddItemHashRange(tableName, hashKey, rangeKey string, item []Attribute) (bool, error) {
	table, _ := GetTable(tableName)

	var attr = make([]dynamodb.Attribute, len(item))
	for i := range item {
		attr[i] = dynamodb.Attribute{
			Type:  item[i].Type,
			Name:  item[i].Name,
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

func buildQueryRange(tableName, hash, operator, rangeValue string) []dynamodb.AttributeComparison {
	schema := GetSchema(tableName)

	var atrrs1 = make([]dynamodb.Attribute, 1)
	atrrs1[0] = dynamodb.Attribute{
		Value: hash,
		Name:  schema.HashKey.Name,
		Type:  schema.HashKey.AttributeType,
	}
	var atrrs2 = make([]dynamodb.Attribute, 1)
	atrrs2[0] = dynamodb.Attribute{
		Value: rangeValue,
		Name:  schema.RangeKey.Name,
		Type:  schema.RangeKey.AttributeType,
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 2)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      schema.HashKey.Name,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs1,
	}
	atrrComaparations[1] = dynamodb.AttributeComparison{
		AttributeName:      schema.RangeKey.Name,
		ComparisonOperator: operator,
		AttributeValueList: atrrs2,
	}

	return atrrComaparations
}

func buildQueryIndexRange(tableName, hash, operator, rangeValue string) []dynamodb.AttributeComparison {
	schema := GetSchema(tableName)

	var atrrs1 = make([]dynamodb.Attribute, 1)
	atrrs1[0] = dynamodb.Attribute{
		Value: hash,
		Name:  schema.GlobalSecondaryIndex.HashKey.Name,
		Type:  schema.GlobalSecondaryIndex.HashKey.AttributeType,
	}
	var atrrs2 = make([]dynamodb.Attribute, 1)
	atrrs2[0] = dynamodb.Attribute{
		Value: rangeValue,
		Name:  schema.GlobalSecondaryIndex.RangeKey.Name,
		Type:  schema.GlobalSecondaryIndex.RangeKey.AttributeType,
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 2)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      schema.GlobalSecondaryIndex.HashKey.Name,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs1,
	}
	atrrComaparations[1] = dynamodb.AttributeComparison{
		AttributeName:      schema.GlobalSecondaryIndex.RangeKey.Name,
		ComparisonOperator: operator,
		AttributeValueList: atrrs2,
	}

	return atrrComaparations
}

func buildQueryRangeBetween(tableName, hash, operator string, rangeValue []string) []dynamodb.AttributeComparison {
	schema := GetSchema(tableName)

	var atrrs1 = make([]dynamodb.Attribute, 1)
	atrrs1[0] = dynamodb.Attribute{
		Value: hash,
		Name:  schema.HashKey.Name,
		Type:  schema.HashKey.AttributeType,
	}
	var atrrs2 = make([]dynamodb.Attribute, 2)
	atrrs2[0] = dynamodb.Attribute{
		Value: rangeValue[0],
		Name:  schema.RangeKey.Name,
		Type:  schema.RangeKey.AttributeType,
	}
	atrrs2[1] = dynamodb.Attribute{
		Value: rangeValue[1],
		Name:  schema.RangeKey.Name,
		Type:  schema.RangeKey.AttributeType,
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 2)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      schema.HashKey.Name,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs1,
	}
	atrrComaparations[1] = dynamodb.AttributeComparison{
		AttributeName:      schema.RangeKey.Name,
		ComparisonOperator: operator,
		AttributeValueList: atrrs2,
	}

	return atrrComaparations
}

func buildQueryIndexRangeBetween(tableName, hash, operator string, rangeValue []string) []dynamodb.AttributeComparison {
	schema := GetSchema(tableName)

	var atrrs1 = make([]dynamodb.Attribute, 1)
	atrrs1[0] = dynamodb.Attribute{
		Value: hash,
		Name:  schema.GlobalSecondaryIndex.HashKey.Name,
		Type:  schema.GlobalSecondaryIndex.HashKey.AttributeType,
	}
	var atrrs2 = make([]dynamodb.Attribute, 2)
	atrrs2[0] = dynamodb.Attribute{
		Value: rangeValue[0],
		Name:  schema.GlobalSecondaryIndex.RangeKey.Name,
		Type:  schema.GlobalSecondaryIndex.RangeKey.AttributeType,
	}
	atrrs2[1] = dynamodb.Attribute{
		Value: rangeValue[1],
		Name:  schema.GlobalSecondaryIndex.RangeKey.Name,
		Type:  schema.GlobalSecondaryIndex.RangeKey.AttributeType,
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 2)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      schema.GlobalSecondaryIndex.HashKey.Name,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs1,
	}
	atrrComaparations[1] = dynamodb.AttributeComparison{
		AttributeName:      schema.GlobalSecondaryIndex.RangeKey.Name,
		ComparisonOperator: operator,
		AttributeValueList: atrrs2,
	}

	return atrrComaparations
}

func buildQueryHash(tableName, hash string) []dynamodb.AttributeComparison {
	schema := GetSchema(tableName)

	var atrrs = make([]dynamodb.Attribute, 1)
	atrrs[0] = dynamodb.Attribute{
		Value: hash,
		Name:  schema.HashKey.Name,
		Type:  schema.HashKey.AttributeType,
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 1)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      schema.HashKey.Name,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs,
	}

	return atrrComaparations
}

func buildQueryIndexHash(tableName, hash string) []dynamodb.AttributeComparison {
	schema := GetSchema(tableName)

	var atrrs = make([]dynamodb.Attribute, 1)
	atrrs[0] = dynamodb.Attribute{
		Value: hash,
		Name:  schema.GlobalSecondaryIndex.HashKey.Name,
		Type:  schema.GlobalSecondaryIndex.HashKey.AttributeType,
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 1)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      schema.GlobalSecondaryIndex.HashKey.Name,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs,
	}

	return atrrComaparations
}

func buildQueryIndexHashRange(tableName, hash string) []dynamodb.AttributeComparison {
	schema := GetSchema(tableName)

	var atrrs = make([]dynamodb.Attribute, 1)
	atrrs[0] = dynamodb.Attribute{
		Value: hash,
		Name:  schema.GlobalSecondaryIndex.HashKey.Name,
		Type:  schema.GlobalSecondaryIndex.HashKey.AttributeType,
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 1)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      schema.GlobalSecondaryIndex.HashKey.Name,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs,
	}

	return atrrComaparations
}

func deleteItem(tableName, hash string) (bool, error) {
	table, _ := GetTable(tableName)

	status, err := table.DeleteItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return status, err
	}

	return status, nil
}

func deleteItems(tableName, hash string, schema Table) (bool, error) {
	table, _ := GetTable(tableName)

	items, _ := RepoGetItemByRange(tableName, hash)
	for i := range items {
		status, err := table.DeleteItem(&dynamodb.Key{HashKey: items[i][schema.HashKey.Name], RangeKey: items[i][schema.RangeKey.Name]})
		if err != nil {
			return status, err
		}
	}

	return true, nil
}
