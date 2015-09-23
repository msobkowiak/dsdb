package main

import (
	"github.com/goamz/goamz/dynamodb"
	"strconv"
)

type DynamoBaseRepository struct {
	table DynamoTable
}

func (r DynamoBaseRepository) GetAll(tableName string) ([]map[string]string, error) {
	table, err := r.table.GetByName(tableName)
	if err != nil {
		return nil, err
	}

	items, err := table.Scan(nil)
	if err != nil {
		return nil, err
	}

	return r.getDataAsArray(items), nil
}

func (r DynamoBaseRepository) GetByHash(tableName, hash string) (map[string]string, error) {
	table, err := r.table.GetByName(tableName)
	if err != nil {
		return nil, err
	}
	item, err := table.GetItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return nil, err
	}

	return r.getData(item), nil
}

func (r DynamoBaseRepository) GetByHashRange(tableName, hashKey, rangeKey string) (map[string]string, error) {
	table, err := r.table.GetByName(tableName)
	if err != nil {
		return nil, err
	}
	item, err := table.GetItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		return nil, err
	}

	return r.getData(item), nil
}

func (r DynamoBaseRepository) GetByOnlyHash(tableName, hashValue string) ([]map[string]string, error) {
	table, err := r.table.GetByName(tableName)
	if err != nil {
		return nil, err
	}
	hashName := table.Key.KeyAttribute.Name

	atrrComaparations := r.buildQueryHash(tableName, hashName, hashValue)

	items, err := table.Query(atrrComaparations)
	if err != nil {
		return nil, err
	}

	return r.getDataAsArray(items), nil
}

func (r DynamoBaseRepository) GetByIndexHash(tableName, indexName, hashValue string) ([]map[string]string, error) {
	table, err := r.table.GetByName(tableName)
	if err != nil {
		return nil, err
	}
	t, _ := schema.GetTableDescription(tableName)
	index, err := t.GetIndexByName(indexName)
	if err != nil {
		return nil, err
	}

	atrrComaparations := r.buildQueryHash(tableName, index.Key.Hash, hashValue)

	items, err := table.QueryOnIndex(atrrComaparations, indexName)
	if err != nil {
		return nil, err
	}

	return r.getDataAsArray(items), nil
}

func (r DynamoBaseRepository) GetByOnlyRange(tableName, hashValue, operator string, rangeValue []string) ([]map[string]string, error) {
	table, err := r.table.GetByName(tableName)
	if err != nil {
		return nil, err
	}
	t, _ := schema.GetTableDescription(tableName)

	var atrrComaparations []dynamodb.AttributeComparison
	atrrComaparations = r.buildQueryRange(tableName, t.PrimaryKey.Hash, hashValue, operator, t.PrimaryKey.Range, rangeValue)
	items, err := table.Query(atrrComaparations)
	if err != nil {
		return nil, err
	}

	return r.getDataAsArray(items), nil
}

func (r DynamoBaseRepository) GetByIndexRange(tableName, indexName, hashValue, operator string, rangeValue []string) ([]map[string]string, error) {
	table, err := r.table.GetByName(tableName)
	if err != nil {
		return nil, err
	}
	t, _ := schema.GetTableDescription(tableName)
	index, err := t.GetIndexByName(indexName)
	if err != nil {
		return nil, err
	}

	var atrrComaparations []dynamodb.AttributeComparison
	atrrComaparations = r.buildQueryRange(tableName, index.Key.Hash, hashValue, operator, index.Key.Range, rangeValue)

	items, err := table.QueryOnIndex(atrrComaparations, indexName)
	if err != nil {
		return nil, err
	}

	return r.getDataAsArray(items), nil
}

func (r DynamoBaseRepository) DeleteByHash(tableName, hash string) (bool, error) {
	t, err := schema.GetTableDescription(tableName)
	if err != nil {
		return false, err
	}

	if t.HasRange() {
		return r.deleteItems(tableName, hash, t)
	}

	return r.deleteItem(tableName, hash)
}

func (r DynamoBaseRepository) DeleteByHashRange(tableName, hashKey, rangeKey string) (bool, error) {
	table, err := r.table.GetByName(tableName)
	if err != nil {
		return false, err
	}

	status, err := table.DeleteItem(&dynamodb.Key{HashKey: hashKey, RangeKey: rangeKey})
	if err != nil {
		return status, err
	}

	return status, nil
}

func (r DynamoBaseRepository) Add(tableName, hashKey, rangeKey string, item []Attribute) (bool, error) {
	t, err := schema.GetTableDescription(tableName)
	if err != nil {
		return false, err
	}
	table, err := r.table.GetByName(tableName)
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
			Value: r.getStringValue(item[i].Value),
		}
	}

	return table.PutItem(hashKey, rangeKey, attr)
}

func (r DynamoBaseRepository) getData(item map[string]*dynamodb.Attribute) map[string]string {
	var data = make(map[string]string)
	for key := range item {
		data[key] = item[key].Value
	}

	return data
}

func (r DynamoBaseRepository) getDataAsArray(items []map[string]*dynamodb.Attribute) []map[string]string {
	itemCount := len(items)
	var data = make([]map[string]string, itemCount)
	for k := 0; k < itemCount; k++ {
		data[k] = r.getData(items[k])
	}

	return data
}

func (r DynamoBaseRepository) buildQueryRange(tableName, hashName, hashValue, operator, rangeName string, rangeValue []string) []dynamodb.AttributeComparison {
	t, _ := schema.GetTableDescription(tableName)
	rangeType := t.GetTypeOfAttribute(rangeName)

	var atrrs1 = make([]dynamodb.Attribute, 1)
	atrrs1[0] = dynamodb.Attribute{
		Value: hashValue,
		Name:  hashName,
		Type:  t.GetTypeOfAttribute(hashName),
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
		Type:  t.GetTypeOfAttribute(rangeName),
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

func (r DynamoBaseRepository) buildQueryHash(tableName, hashName, hashValue string) []dynamodb.AttributeComparison {
	t, _ := schema.GetTableDescription(tableName)

	var atrrs = make([]dynamodb.Attribute, 1)
	atrrs[0] = dynamodb.Attribute{
		Value: hashValue,
		Name:  hashName,
		Type:  t.GetTypeOfAttribute(hashName),
	}

	var atrrComaparations = make([]dynamodb.AttributeComparison, 1)
	atrrComaparations[0] = dynamodb.AttributeComparison{
		AttributeName:      hashName,
		ComparisonOperator: "EQ",
		AttributeValueList: atrrs,
	}
	return atrrComaparations
}

func (r DynamoBaseRepository) deleteItem(tableName, hash string) (bool, error) {
	table, _ := r.table.GetByName(tableName)

	status, err := table.DeleteItem(&dynamodb.Key{HashKey: hash})
	if err != nil {
		return status, err
	}

	return status, nil
}

func (r DynamoBaseRepository) deleteItems(tableName, hash string, schema TableDescription) (bool, error) {
	table, _ := r.table.GetByName(tableName)

	items, _ := r.GetByOnlyHash(tableName, hash)
	for i := range items {
		status, err := table.DeleteItem(&dynamodb.Key{HashKey: items[i][schema.PrimaryKey.Hash], RangeKey: items[i][schema.PrimaryKey.Range]})
		if err != nil {
			return status, err
		}
	}

	return true, nil
}

func (r DynamoBaseRepository) getStringValue(itemValue interface{}) string {
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
