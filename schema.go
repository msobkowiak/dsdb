package main

type Table struct {
	Name               string
	HashKey            Key
	RangeKey           Key
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

type Key struct {
	Name          string
	AttributeType string
	KeyType       string
	Value         string
}

func (t Table) HasRange() bool {
	if t.RangeKey.Name != "" {
		return true
	}

	return false
}

func (t Table) GetRangeName() string {
	return t.HashKey.Name
}
