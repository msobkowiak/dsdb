package main

import (
	"github.com/goamz/goamz/dynamodb"
)

var tables = map[string]TableDescription{
	"users": TableDescription{
		Name: "users",
		Attributes: []AttributeDefinition{
			AttributeDefinition{"id", "N", true},
			AttributeDefinition{"email", "S", true},
		},
		PrimaryKey: PrimaryKeyDefinition{
			Type: "HASH",
			Hash: "id",
		},
		SecondaryIndexes: []SecondaryIndexDefinition{
			SecondaryIndexDefinition{
				Name: "email",
				Type: "HASH",
				Hash: "email",
			},
		},
		Authentication: Authentication{
			DynamoAuth{
				Region:    "http://127.0.0.1:4567",
				AccessKey: "access",
				SecretKey: "secret",
			},
		},
	},
	"game_scores": TableDescription{
		Name: "game_scores",
		Attributes: []AttributeDefinition{
			AttributeDefinition{"user_id", "N", true},
			AttributeDefinition{"game_title", "S", true},
			AttributeDefinition{"wins", "N", true},
			AttributeDefinition{"losts", "N", true},
		},
		PrimaryKey: PrimaryKeyDefinition{
			Type:  "RANGE",
			Hash:  "game_title",
			Range: "user_id",
		},
		SecondaryIndexes: []SecondaryIndexDefinition{
			SecondaryIndexDefinition{
				Name:  "wins_losts",
				Type:  "RANGE",
				Hash:  "wins",
				Range: "losts",
			},
		},
		Authentication: Authentication{
			DynamoAuth{
				Region:    "http://127.0.0.1:4567",
				AccessKey: "access",
				SecretKey: "secret",
			},
		},
	},
}

func GetTableDescription(tableName string) TableDescription {
	return tables[tableName]
}

func GetTableData(tableName string) ([][]dynamodb.Attribute, []string) {
	return data[tableName], hashKeys[tableName]
}

var data = map[string][][]dynamodb.Attribute{
	"users": [][]dynamodb.Attribute{
		[]dynamodb.Attribute{
			*dynamodb.NewStringAttribute("first_name", "Monika"),
			*dynamodb.NewStringAttribute("last_name", "Sobkowiak"),
			*dynamodb.NewStringAttribute("email", "monika@gmail.com"),
			*dynamodb.NewStringAttribute("country", "Poland"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewStringAttribute("first_name", "Ana"),
			*dynamodb.NewStringAttribute("last_name", "Dias"),
			*dynamodb.NewStringAttribute("email", "ana@gmail.com"),
			*dynamodb.NewStringAttribute("country", "Portugal"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewStringAttribute("first_name", "Nuno"),
			*dynamodb.NewStringAttribute("last_name", "Correia"),
			*dynamodb.NewStringAttribute("email", "nuno@exemple.com"),
			*dynamodb.NewStringAttribute("country", "Portugal"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewStringAttribute("first_name", "Isabel"),
			*dynamodb.NewStringAttribute("last_name", "Frenandes"),
			*dynamodb.NewStringAttribute("email", "isabel@gmail.com"),
			*dynamodb.NewStringAttribute("country", "Spain"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewStringAttribute("first_name", "Miguel"),
			*dynamodb.NewStringAttribute("last_name", "Oliveira"),
			*dynamodb.NewStringAttribute("email", "miguel@gmail.com"),
			*dynamodb.NewStringAttribute("country", "Portugal"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewStringAttribute("first_name", "Mikolaj"),
			*dynamodb.NewStringAttribute("last_name", "Nowak"),
			*dynamodb.NewStringAttribute("email", "mikolaj@exemple.com"),
			*dynamodb.NewStringAttribute("country", "Poland"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewStringAttribute("first_name", "Joao"),
			*dynamodb.NewStringAttribute("last_name", "Silva"),
			*dynamodb.NewStringAttribute("email", "joao@gmail.com"),
			*dynamodb.NewStringAttribute("country", "Portugal"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewStringAttribute("first_name", "Mat"),
			*dynamodb.NewStringAttribute("last_name", "Deamon"),
			*dynamodb.NewStringAttribute("email", "mat@gmail.com"),
			*dynamodb.NewStringAttribute("country", "USA"),
		},
	},
	"game_scores": [][]dynamodb.Attribute{
		[]dynamodb.Attribute{
			*dynamodb.NewNumericAttribute("top_score", "5842"),
			*dynamodb.NewNumericAttribute("wins", "8"),
			*dynamodb.NewNumericAttribute("losts", "2"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewNumericAttribute("top_score", "123"),
			*dynamodb.NewNumericAttribute("wins", "3"),
			*dynamodb.NewNumericAttribute("losts", "0"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewNumericAttribute("top_score", "333333"),
			*dynamodb.NewNumericAttribute("wins", "30"),
			*dynamodb.NewNumericAttribute("losts", "90"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewNumericAttribute("top_score", "12"),
			*dynamodb.NewNumericAttribute("wins", "2"),
			*dynamodb.NewNumericAttribute("losts", "2"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewNumericAttribute("top_score", "45"),
			*dynamodb.NewNumericAttribute("wins", "5"),
			*dynamodb.NewNumericAttribute("losts", "1"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewNumericAttribute("top_score", "667854"),
			*dynamodb.NewNumericAttribute("wins", "399"),
			*dynamodb.NewNumericAttribute("losts", "100"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewNumericAttribute("top_score", "23"),
			*dynamodb.NewNumericAttribute("wins", "1"),
			*dynamodb.NewNumericAttribute("losts", "30"),
		},
		[]dynamodb.Attribute{
			*dynamodb.NewNumericAttribute("top_score", "58542"),
			*dynamodb.NewNumericAttribute("wins", "70"),
			*dynamodb.NewNumericAttribute("losts", "2"),
		},
	},
}

var hashKeys = map[string][]string{
	"game_scores": []string{
		"Mario Brodes",
		"Medal of Honor",
		"Game X",
		"Mario Brodes",
		"Game X",
		"Game Y",
		"Game Y",
		"Mario Brodes",
	},
}
