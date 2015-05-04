package main

type Table struct {
	Name                 string
	HashKey              Key
	RangeKey             Key
	GlobalSecondaryIndex SecondaryIndex
	ReadCapacityUnits    int64
	WriteCapacityUnits   int64
}

type Key struct {
	Name          string
	AttributeType string
	KeyType       string
	Value         string
}

type Attribute struct {
	Type  string
	Name  string
	Value string
}

type SecondaryIndex struct {
	Name     string
	HashKey  Key
	RangeKey Key
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
