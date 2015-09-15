package main

import (
	"github.com/goamz/goamz/dynamodb"
	"strconv"
)

type DynamoCRUDRepository struct {
}

func (r DynamoCRUDRepository) GetAll(tableName string) ([]map[string]string, error) {
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return nil, err
	}

	items, err := table.Scan(nil)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func (r DynamoCRUDRepository) GetByHash(tableName, hash string) (map[string]string, error) {
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return nil, err
	}
	item, err := table.GetItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return nil, err
	}

	return getData(item), nil
}

func (r DynamoCRUDRepository) GetByHashRange(tableName, hashKey, rangeKey string) (map[string]string, error) {
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return nil, err
	}
	item, err := table.GetItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		return nil, err
	}

	return getData(item), nil
}

func (r DynamoCRUDRepository) GetByOnlyHash(tableName, hashValue string) ([]map[string]string, error) {
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return nil, err
	}
	hashName := table.Key.KeyAttribute.Name

	atrrComaparations := buildQueryHash(tableName, hashName, hashValue)

	items, err := table.Query(atrrComaparations)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func (r DynamoCRUDRepository) GetByIndexHash(tableName, indexName, hashValue string) ([]map[string]string, error) {
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return nil, err
	}
	schema, _ := GetTableDescription(tableName, schema.Tables)
	index, err := schema.GetIndexByName(indexName)
	if err != nil {
		return nil, err
	}

	atrrComaparations := buildQueryHash(tableName, index.Key.Hash, hashValue)

	items, err := table.QueryOnIndex(atrrComaparations, indexName)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func (r DynamoCRUDRepository) GetByOnlyRange(tableName, hashValue, operator string, rangeValue []string) ([]map[string]string, error) {
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return nil, err
	}
	schema, _ := GetTableDescription(tableName, schema.Tables)

	var atrrComaparations []dynamodb.AttributeComparison
	atrrComaparations = buildQueryRange(tableName, schema.PrimaryKey.Hash, hashValue, operator, schema.PrimaryKey.Range, rangeValue)
	items, err := table.Query(atrrComaparations)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func (r DynamoCRUDRepository) GetByIndexRange(tableName, indexName, hashValue, operator string, rangeValue []string) ([]map[string]string, error) {
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return nil, err
	}
	schema, _ := GetTableDescription(tableName, schema.Tables)
	index, err := schema.GetIndexByName(indexName)
	if err != nil {
		return nil, err
	}

	var atrrComaparations []dynamodb.AttributeComparison
	atrrComaparations = buildQueryRange(tableName, index.Key.Hash, hashValue, operator, index.Key.Range, rangeValue)

	items, err := table.QueryOnIndex(atrrComaparations, indexName)
	if err != nil {
		return nil, err
	}

	return getDataAsArray(items), nil
}

func (r DynamoCRUDRepository) DeleteByHash(tableName, hash string) (bool, error) {
	schema, err := GetTableDescription(tableName, schema.Tables)
	if err != nil {
		return false, err
	}

	if schema.HasRange() {
		return deleteItems(tableName, hash, schema)
	}

	return deleteItem(tableName, hash)
}

func (r DynamoCRUDRepository) DeleteByHashRange(tableName, hashKey, rangeKey string) (bool, error) {
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return false, err
	}

	status, err := table.DeleteItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		return status, err
	}

	return status, nil
}

func (r DynamoCRUDRepository) Add(tableName, hashKey, rangeKey string, item []Attribute) (bool, error) {
	t, err := GetTableDescription(tableName, schema.Tables)
	if err != nil {
		return false, err
	}
	table, err := GetDynamoTable(tableName)
	if err != nil {
		return false, err
	}

	var attr = make([]dynamodb.Attribute, len(item))
	for i := range item {
		name := item[i].Description.Name
		var attrType string
		if t.GetTypeOfAttribute(name) == "G" {
			attrType = "S"
		} else {
			attrType = t.GetTypeOfAttribute(name)
		}
		attr[i] = dynamodb.Attribute{
			Type:  attrType,
			Name:  item[i].Description.Name,
			Value: getStringValue(item[i].Value),
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
	schema, _ := GetTableDescription(tableName, schema.Tables)
	rangeType := schema.GetTypeOfAttribute(rangeName)

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
			Type:  rangeType,
		}
	} else if len(rangeValue) == 2 {
		atrrs2 = make([]dynamodb.Attribute, 2)
		atrrs2[0] = dynamodb.Attribute{
			Value: rangeValue[0],
			Name:  rangeName,
			Type:  rangeType,
		}
		atrrs2[1] = dynamodb.Attribute{
			Value: rangeValue[1],
			Name:  rangeName,
			Type:  rangeType,
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
	schema, _ := GetTableDescription(tableName, schema.Tables)

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
	table, _ := GetDynamoTable(tableName)

	status, err := table.DeleteItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return status, err
	}

	return status, nil
}

func deleteItems(tableName, hash string, schema TableDescription) (bool, error) {
	table, _ := GetDynamoTable(tableName)
	var repo DynamoCRUDRepository

	items, _ := repo.GetByOnlyHash(tableName, hash)
	for i := range items {
		status, err := table.DeleteItem(&dynamodb.Key{HashKey: items[i][schema.PrimaryKey.Hash], RangeKey: items[i][schema.PrimaryKey.Range]})
		if err != nil {
			return status, err
		}
	}

	return true, nil
}

func getStringValue(itemValue interface{}) string {
	var value string

	switch itemValue.(type) {
	case string:
		value = itemValue.(string)
	case int, int32, int64:
		value = strconv.FormatInt(int64(itemValue.(int)), 10)
	case float64:
		value = strconv.FormatFloat(itemValue.(float64), 'f', 8, 64)
	}

	return value
}
