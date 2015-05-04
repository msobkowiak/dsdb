package main

import (
	"github.com/goamz/goamz/dynamodb"
)

func GetSchema(tableName string) Table {
	if tableName == "users" {
		return getUsersSchema()
	} else if tableName == "game_scores" {
		return getGameScoreSchema()
	}

	return Table{
		Name: tableName,
	}
}

func getUsersSchema() Table {
	var t Table
	t.Name = "users"
	t.HashKey.Name = "id"
	t.HashKey.AttributeType = "N"
	t.HashKey.KeyType = "HASH"
	t.ReadCapacityUnits = 10
	t.WriteCapacityUnits = 10

	return t
}

func getGameScoreSchema() Table {
	var t Table
	t.Name = "game_scores"
	t.HashKey.Name = "game_title"
	t.HashKey.AttributeType = "S"
	t.HashKey.KeyType = "HASH"
	t.RangeKey.Name = "user_id"
	t.RangeKey.AttributeType = "N"
	t.RangeKey.KeyType = "RANGE"
	t.ReadCapacityUnits = 10
	t.WriteCapacityUnits = 10
	t.GlobalSecondaryIndex.Name = "wins"
	t.GlobalSecondaryIndex.HashKey.Name = "game_title"
	t.GlobalSecondaryIndex.HashKey.AttributeType = "S"
	t.GlobalSecondaryIndex.HashKey.KeyType = "HASH"
	t.GlobalSecondaryIndex.RangeKey.Name = "wins"
	t.GlobalSecondaryIndex.RangeKey.AttributeType = "N"
	t.GlobalSecondaryIndex.RangeKey.KeyType = "RANGE"
	return t
}

func LoadUsersData() [][]dynamodb.Attribute {
	var data = make([][]dynamodb.Attribute, 8)
	data[0] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Monika"),
		*dynamodb.NewStringAttribute("last_name", "Sobkowiak"),
		*dynamodb.NewStringAttribute("email", "monika@gmail.com"),
		*dynamodb.NewStringAttribute("country", "Poland"),
	}
	data[1] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Ana"),
		*dynamodb.NewStringAttribute("last_name", "Dias"),
		*dynamodb.NewStringAttribute("email", "ana@gmail.com"),
		*dynamodb.NewStringAttribute("country", "Portugal"),
	}
	data[2] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Nuno"),
		*dynamodb.NewStringAttribute("last_name", "Correia"),
		*dynamodb.NewStringAttribute("email", "nuno@exemple.com"),
		*dynamodb.NewStringAttribute("country", "Portugal"),
	}
	data[3] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Isabel"),
		*dynamodb.NewStringAttribute("last_name", "Frenandes"),
		*dynamodb.NewStringAttribute("email", "isabel@gmail.com"),
		*dynamodb.NewStringAttribute("country", "Spain"),
	}
	data[4] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Miguel"),
		*dynamodb.NewStringAttribute("last_name", "Oliveira"),
		*dynamodb.NewStringAttribute("email", "miguel@gmail.com"),
		*dynamodb.NewStringAttribute("country", "Portugal"),
	}
	data[5] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Mikolaj"),
		*dynamodb.NewStringAttribute("last_name", "Nowak"),
		*dynamodb.NewStringAttribute("email", "mikolaj@exemple.com"),
		*dynamodb.NewStringAttribute("country", "Poland"),
	}
	data[6] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Joao"),
		*dynamodb.NewStringAttribute("last_name", "Silva"),
		*dynamodb.NewStringAttribute("email", "joao@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Portugal"),
	}
	data[7] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Mat"),
		*dynamodb.NewStringAttribute("last_name", "Deamon"),
		*dynamodb.NewStringAttribute("email", "mat@gmail.com"),
		*dynamodb.NewStringAttribute("country", "USA"),
	}

	return data
}

func LoadGameScoreData() ([][]dynamodb.Attribute, []string) {
	var data = make([][]dynamodb.Attribute, 8)
	data[0] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "5842"),
		*dynamodb.NewNumericAttribute("wins", "8"),
		*dynamodb.NewNumericAttribute("losts", "2"),
	}
	data[1] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "123"),
		*dynamodb.NewNumericAttribute("wins", "3"),
		*dynamodb.NewNumericAttribute("losts", "0"),
	}
	data[2] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "333333"),
		*dynamodb.NewNumericAttribute("wins", "30"),
		*dynamodb.NewNumericAttribute("losts", "90"),
	}
	data[3] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "12"),
		*dynamodb.NewNumericAttribute("wins", "2"),
		*dynamodb.NewNumericAttribute("losts", "2"),
	}
	data[4] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "45"),
		*dynamodb.NewNumericAttribute("wins", "5"),
		*dynamodb.NewNumericAttribute("losts", "1"),
	}
	data[5] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "667854"),
		*dynamodb.NewNumericAttribute("wins", "399"),
		*dynamodb.NewNumericAttribute("losts", "100"),
	}
	data[6] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "23"),
		*dynamodb.NewNumericAttribute("wins", "1"),
		*dynamodb.NewNumericAttribute("losts", "30"),
	}
	data[7] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "58542"),
		*dynamodb.NewNumericAttribute("wins", "70"),
		*dynamodb.NewNumericAttribute("losts", "2"),
	}

	var hashKeys = make([]string, 8)
	hashKeys[0] = "Mario Brodes"
	hashKeys[1] = "Medal of Honor"
	hashKeys[2] = "Game X"
	hashKeys[3] = "Mario Brodes"
	hashKeys[4] = "Game X"
	hashKeys[5] = "Game Y"
	hashKeys[6] = "Game Y"
	hashKeys[7] = "Mario Brodes"

	return data, hashKeys
}
