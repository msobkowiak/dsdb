package main

import (
	"errors"
)

const defaultDataType = "S"

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

	return defaultDataType
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

func (s DbDescription) GetTableDescription(tableName string) (TableDescription, error) {
	if s.Tables[tableName].Name != "" {
		return s.Tables[tableName], nil
	} else {
		return TableDescription{}, errors.New("Table " + tableName + " not found.")
	}
}

func (s DbDescription) GetHashName(tableName string) (string, error) {
	table, err := s.GetTableDescription(tableName)
	if err != nil {
		return "", err
	}

	return table.PrimaryKey.Hash, nil
}

func (s DbDescription) GetRangeName(tableName string) (string, error) {
	table, err := s.GetTableDescription(tableName)
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

func (t TableDescription) ExcludeNonKeyAttributes() []AttributeDefinition {
	var newAttr = []AttributeDefinition{}
	for _, attr := range t.Attributes {
		if t.isKeySchemaAttribute(attr) {
			newAttr = append(newAttr, attr)
		}
	}

	return newAttr
}

func (t TableDescription) isKeySchemaAttribute(attr AttributeDefinition) bool {
	if t.isPrimaryKeyAttribute(attr) || t.isSecondaryIndexAttribute(attr) {
		return true
	}

	return false
}

func (t TableDescription) isPrimaryKeyAttribute(attr AttributeDefinition) bool {
	if attr.Name == t.PrimaryKey.Hash || attr.Name == t.PrimaryKey.Range {
		return true
	}

	return false
}

func (t TableDescription) isSecondaryIndexAttribute(attr AttributeDefinition) bool {
	for _, index := range t.SecondaryIndexes {
		if attr.Name == index.Key.Hash || attr.Name == index.Key.Range {
			return true
		}
	}

	return false
}
