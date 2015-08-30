package main

import (
	"errors"
)

type DbDescription struct {
	Name           string
	Authentication Authentication
	Tables         map[string]TableDescription
}

type TableDescription struct {
	Name             string
	Attributes       []AttributeDefinition
	PrimaryKey       KeyDefinition
	SecondaryIndexes []SecondaryIndexDefinition
}

type AttributeDefinition struct {
	Name     string
	Type     string
	Required bool
}

type KeyDefinition struct {
	Type  string
	Hash  string
	Range string
}

type SecondaryIndexDefinition struct {
	Name string
	Key  KeyDefinition
}

type Attribute struct {
	Description AttributeDefinition
	Value       interface{}
}

type Authentication struct {
	Dynamo DynamoAuth
}

type DynamoAuth struct {
	Region    string
	AccessKey string
	SecretKey string
}

type TextSearch struct {
	Field     string
	Query     string
	Operator  string
	Precision string
}

func (t TableDescription) HasRange() bool {
	if t.PrimaryKey.Range != "" {
		return true
	}

	return false
}

func (t TableDescription) GetTypeOfAttribute(name string) string {
	for i := range t.Attributes {
		if t.Attributes[i].Name == name {
			return t.Attributes[i].Type
		}
	}

	// Return the defualt value
	return "S"
}

func (t TableDescription) GetIndexByName(name string) (SecondaryIndexDefinition, error) {
	for i := range t.SecondaryIndexes {
		if t.SecondaryIndexes[i].Name == name {
			return t.SecondaryIndexes[i], nil
		}
	}

	return SecondaryIndexDefinition{}, errors.New("Index not found")
}

func (t TableDescription) HasGeoPoint() bool {
	for _, attr := range t.Attributes {
		if attr.Type == "G" {
			return true
		}
	}

	return false
}

func GetTableDescription(tableName string, tables map[string]TableDescription) (TableDescription, error) {
	if tables[tableName].Name != "" {
		return tables[tableName], nil
	} else {
		return TableDescription{}, errors.New("Table " + tableName + " not found.")
	}
}

func GetHashName(tableName string, schema DbDescription) (string, error) {
	table, err := GetTableDescription(tableName, schema.Tables)
	if err != nil {
		return "", err
	}

	return table.PrimaryKey.Hash, nil
}

func GetRangeName(tableName string, schema DbDescription) (string, error) {
	table, err := GetTableDescription(tableName, schema.Tables)
	if err != nil {
		return "", err
	}

	return table.PrimaryKey.Range, nil
}

func (t TableDescription) GetGeoPointName() (string, error) {
	for _, atrr := range t.Attributes {
		if atrr.Type == "G" {
			return atrr.Name, nil
		}
	}

	return "", errors.New("No geo point found")
}

func ExcludeNonKeyAttributes(table TableDescription) []AttributeDefinition {
	var newAttr = []AttributeDefinition{}
	for _, attr := range table.Attributes {
		if isKeySchemaAttribute(attr, table) {
			newAttr = append(newAttr, attr)
		}
	}

	return newAttr
}

func isKeySchemaAttribute(attr AttributeDefinition, table TableDescription) bool {
	if isPrimaryKeyAttribute(attr, table) || isSecondaryIndexAttribute(attr, table) {
		return true
	}

	return false
}

func isPrimaryKeyAttribute(attr AttributeDefinition, table TableDescription) bool {
	if attr.Name == table.PrimaryKey.Hash || attr.Name == table.PrimaryKey.Range {
		return true
	}

	return false
}

func isSecondaryIndexAttribute(attr AttributeDefinition, table TableDescription) bool {
	for _, index := range table.SecondaryIndexes {
		if attr.Name == index.Key.Hash || attr.Name == index.Key.Range {
			return true
		}
	}

	return false
}
