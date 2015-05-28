package main

import (
	"strconv"
)

var dbDescription = DbDescription{
	Name: "test",
	Authentication: Authentication{
		DynamoAuth{
			Region:    "http://127.0.0.1:4567",
			AccessKey: "access",
			SecretKey: "secret",
		},
	},
	Tables: map[string]TableDescription{
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
		},
	},
}

var data = map[string][][]Attribute{
	"users": [][]Attribute{
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Monika"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Sobkowiak"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "monika@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Poland"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Ana"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Dias"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "ana@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Nuno"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Correia"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "nuno@example.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Nuno"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Correia"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "nuno@example.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Isabel"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Fernendes"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "isabel@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Spain"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Miguel"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Oliveira"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "miguel@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Mikolaj"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Nowak"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "mikolaj@example.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Poland"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Joao"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Silva"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "joao@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
	},
	"game_scores": [][]Attribute{
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: "123"},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: "20"},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: "0"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: "333"},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: "90"},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: "21"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: "444"},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: "99"},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: "59"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: "555"},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: "12"},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: "9"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: "666"},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: "7"},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: "20"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: "122"},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: "20"},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: "7"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: "777"},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: "190"},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: "87"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: "888"},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: "438"},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: "164"},
		},
	},
}

var hashKeys = map[string][]string{
	"game_scores": []string{
		"Game X",
		"Game Y",
		"Game X",
		"Game Z",
		"Game X",
		"Game Y",
		"Game Y",
		"Game Z",
	},
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type data1 struct {
	Database_name     string
	Dynamo_region     string
	Dynamo_access_key string
	Dynamo_secret_key string
	Tables            []struct {
		Name              string
		Attributes        []AttributeDefinition
		Primary_key       PrimaryKeyDefinition
		Secondary_indexes []SecondaryIndexDefinition
	}
}

func Bootstrap() {

	db := Auth("http://127.0.0.1:4567", "key", "secret")

	// cleanup the database
	DeleteAllTables(db)

	//create tables with example data
	tableName := "users"
	CreateTable(dbDescription.Tables[tableName])
	for i := range data[tableName] {
		hash := strconv.FormatInt(int64(i+1), 10)
		RepoAddItem(tableName, hash, data[tableName][i])
		AddToElasticSearch(dbDescription.Name, tableName, hash, "", data[tableName][i])
	}

	tableName = "game_scores"
	CreateTable(dbDescription.Tables[tableName])
	for i := range data[tableName] {
		rangeValue := strconv.FormatInt(int64(i+1), 10)
		RepoAddItemHashRange(tableName, hashKeys[tableName][i], rangeValue, data[tableName][i])
		AddToElasticSearch(dbDescription.Name, tableName, hashKeys[tableName][i], rangeValue, data[tableName][i])
	}
}

func LoadSchema() DbDescription {
	return dbDescription
}
