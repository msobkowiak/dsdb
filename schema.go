package main

type Table struct {
	Name               string
	PrimaryKey         Key
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

type Key struct {
	Name          string
	AttributeType string
	KeyType       string
	Value         string
}
