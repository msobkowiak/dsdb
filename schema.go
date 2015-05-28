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
	PrimaryKey       PrimaryKeyDefinition
	SecondaryIndexes []SecondaryIndexDefinition
}

type AttributeDefinition struct {
	Name     string
	Type     string
	Required bool
}

type PrimaryKeyDefinition struct {
	Type  string
	Hash  string
	Range string
}

type SecondaryIndexDefinition struct {
	Name  string
	Type  string
	Hash  string
	Range string
}

type Attribute struct {
	Description AttributeDefinition
	Value       string
}

type Authentication struct {
	Dynamo DynamoAuth
}

type DynamoAuth struct {
	Region    string
	AccessKey string
	SecretKey string
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

	return ""
}

func (t TableDescription) GetIndexByName(name string) (SecondaryIndexDefinition, error) {
	for i := range t.SecondaryIndexes {
		if t.SecondaryIndexes[i].Name == name {
			return t.SecondaryIndexes[i], nil
		}
	}

	return SecondaryIndexDefinition{}, errors.New("Index not found")
}

func GetTableDescription(tableName string, tables map[string]TableDescription) (TableDescription, error) {
	if tables[tableName].Name != "" {
		return tables[tableName], nil
	} else {
		return TableDescription{}, errors.New("Table " + tableName + " not found.")
	}
}
