package main

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
