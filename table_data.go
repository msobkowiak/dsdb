package main

import (
	"github.com/goamz/goamz/dynamodb"
)

func GetSchema(tableName string) Table {
	if tableName == "users" {
		return getUsersSchema()
	}

	return getGameScoreSchema()
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
	t.HashKey.Name = "user_id"
	t.HashKey.AttributeType = "N"
	t.HashKey.KeyType = "HASH"
	t.RangeKey.Name = "game_title"
	t.RangeKey.AttributeType = "S"
	t.RangeKey.KeyType = "RANGE"
	t.ReadCapacityUnits = 10
	t.WriteCapacityUnits = 10

	return t
}

func LoadUsersData() [][]dynamodb.Attribute {
	var data = make([][]dynamodb.Attribute, 8)
	data[0] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Monika"),
		*dynamodb.NewStringAttribute("last_name", "Sobkowiak"),
		*dynamodb.NewStringAttribute("email", "monika@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Poland"),
	}
	data[1] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Ana"),
		*dynamodb.NewStringAttribute("last_name", "Dias"),
		*dynamodb.NewStringAttribute("email", "ana@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Portugal"),
	}
	data[2] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Nuno"),
		*dynamodb.NewStringAttribute("last_name", "Correia"),
		*dynamodb.NewStringAttribute("email", "nuno@exemple.com"),
		*dynamodb.NewStringAttribute("counrty", "Portugal"),
	}
	data[3] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Isabel"),
		*dynamodb.NewStringAttribute("last_name", "Frenandes"),
		*dynamodb.NewStringAttribute("email", "isabel@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Spain"),
	}
	data[4] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Miguel"),
		*dynamodb.NewStringAttribute("last_name", "Oliveira"),
		*dynamodb.NewStringAttribute("email", "miguel@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Portugal"),
	}
	data[5] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Mikolaj"),
		*dynamodb.NewStringAttribute("last_name", "Nowak"),
		*dynamodb.NewStringAttribute("email", "mikolaj@exemple.com"),
		*dynamodb.NewStringAttribute("counrty", "Poland"),
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
		*dynamodb.NewStringAttribute("counrty", "USA"),
	}

	return data
}

func LoadGameScoreData() ([][]dynamodb.Attribute, []string) {
	var data = make([][]dynamodb.Attribute, 8)
	data[0] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "5842"),
		*dynamodb.NewNumericAttribute("wins", "8"),
	}
	data[1] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "123"),
		*dynamodb.NewNumericAttribute("wins", "3"),
	}
	data[2] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "333333"),
		*dynamodb.NewNumericAttribute("wins", "30"),
	}
	data[3] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "12"),
		*dynamodb.NewNumericAttribute("wins", "2"),
	}
	data[4] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "45"),
		*dynamodb.NewNumericAttribute("wins", "5"),
	}
	data[5] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "667854"),
		*dynamodb.NewNumericAttribute("wins", "399"),
	}
	data[6] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "23"),
		*dynamodb.NewNumericAttribute("wins", "1"),
	}
	data[7] = []dynamodb.Attribute{
		*dynamodb.NewNumericAttribute("top_score", "58542"),
		*dynamodb.NewNumericAttribute("wins", "70"),
	}

	var rangeKeys = make([]string, 8)
	rangeKeys[0] = "Mario Brodes"
	rangeKeys[1] = "Medal of Honor"
	rangeKeys[2] = "Game X"
	rangeKeys[3] = "Mario Brodes"
	rangeKeys[4] = "Game X"
	rangeKeys[5] = "Game Y"
	rangeKeys[6] = "Game Y"
	rangeKeys[7] = "Mario Brodes"

	return data, rangeKeys
}
