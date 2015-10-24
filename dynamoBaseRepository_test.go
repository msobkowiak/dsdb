package main

import (
	"github.com/goamz/goamz/dynamodb"
	. "gopkg.in/check.v1"
)

type DynamoBaseRepositorySuite struct {
	Db       DbDescription
	Data     map[string][][]dynamodb.Attribute
	HashKeys map[string][]string
}

func (s *DynamoBaseRepositorySuite) SetUpSuite(c *C) {
	schema = dynamoBaseRepository_suite.Db
	var t DynamoTable

	t.Create(schema.Tables["test"])
	users := t.Create(schema.Tables["users_test"])
	games := t.Create(schema.Tables["game_scores_test"])
	t.AddItems(users, dynamoBaseRepository_suite.Data["users_test"], nil)
	t.AddItems(games, dynamoBaseRepository_suite.Data["game_scores_test"], dynamoBaseRepository_suite.HashKeys["game_scores_test"])
}

func (s *DynamoBaseRepositorySuite) TearDownSuite(c *C) {
	var client DynamoClient
	dynamAuth := dynamoBaseRepository_suite.Db.Authentication.Dynamo
	db := client.Auth(dynamAuth.Region, dynamAuth.AccessKey, dynamAuth.SecretKey)
	client.DeleteAll(db)
}

var dynamoBaseRepository_suite = &DynamoBaseRepositorySuite{

	Db: DbDescription{
		Name: "test",
		Authentication: Authentication{
			DynamoAuth{
				Region:    "http://127.0.0.1:4567",
				AccessKey: "access",
				SecretKey: "secret",
			},
		},
		Tables: map[string]TableDescription{
			"users_test": TableDescription{
				Name: "users_test",
				Attributes: []AttributeDefinition{
					AttributeDefinition{"id", "N", true},
					AttributeDefinition{"email", "S", true},
					AttributeDefinition{"country", "S", true},
				},
				PrimaryKey: KeyDefinition{
					Type: "HASH",
					Hash: "id",
				},
				SecondaryIndexes: []SecondaryIndexDefinition{
					SecondaryIndexDefinition{
						Name: "email",
						Key: KeyDefinition{
							Type: "HASH",
							Hash: "email",
						},
					},
					SecondaryIndexDefinition{
						Name: "country",
						Key: KeyDefinition{
							Type: "HASH",
							Hash: "country",
						},
					},
				},
			},
			"game_scores_test": TableDescription{
				Name: "game_scores_test",
				Attributes: []AttributeDefinition{
					AttributeDefinition{"user_id", "N", true},
					AttributeDefinition{"game_title", "S", true},
					AttributeDefinition{"wins", "N", true},
					AttributeDefinition{"losts", "N", true},
				},
				PrimaryKey: KeyDefinition{
					Type:  "RANGE",
					Hash:  "game_title",
					Range: "user_id",
				},
				SecondaryIndexes: []SecondaryIndexDefinition{
					SecondaryIndexDefinition{
						Name: "wins_losts",
						Key: KeyDefinition{
							Type:  "RANGE",
							Hash:  "wins",
							Range: "losts",
						},
					},
				},
			},
			"test": TableDescription{
				Name: "test",
				Attributes: []AttributeDefinition{
					AttributeDefinition{"id", "N", true},
					AttributeDefinition{"name", "S", true},
				},
				PrimaryKey: KeyDefinition{
					Type: "HASH",
					Hash: "id",
				},
			},
		},
	},
	Data: map[string][][]dynamodb.Attribute{
		"users_test": [][]dynamodb.Attribute{
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
		},
		"game_scores_test": [][]dynamodb.Attribute{
			[]dynamodb.Attribute{
				*dynamodb.NewNumericAttribute("top_score", "582"),
				*dynamodb.NewNumericAttribute("wins", "8"),
				*dynamodb.NewNumericAttribute("losts", "2"),
			},
			[]dynamodb.Attribute{
				*dynamodb.NewNumericAttribute("top_score", "123"),
				*dynamodb.NewNumericAttribute("wins", "8"),
				*dynamodb.NewNumericAttribute("losts", "0"),
			},
			[]dynamodb.Attribute{
				*dynamodb.NewNumericAttribute("top_score", "333333"),
				*dynamodb.NewNumericAttribute("wins", "8"),
				*dynamodb.NewNumericAttribute("losts", "90"),
			},
			[]dynamodb.Attribute{
				*dynamodb.NewNumericAttribute("top_score", "12"),
				*dynamodb.NewNumericAttribute("wins", "8"),
				*dynamodb.NewNumericAttribute("losts", "2"),
			},
		},
		"geo_test": [][]dynamodb.Attribute{
			[]dynamodb.Attribute{
				*dynamodb.NewNumericAttribute("name", "place1"),
				*dynamodb.NewNumericAttribute("description", "place1 description"),
				*dynamodb.NewNumericAttribute("location", "41.158915,-8.631091"),
			},
			[]dynamodb.Attribute{
				*dynamodb.NewNumericAttribute("name", "place2"),
				*dynamodb.NewNumericAttribute("description", "place2 description"),
				*dynamodb.NewNumericAttribute("location", "42.158915,-7.631091"),
			},
			[]dynamodb.Attribute{
				*dynamodb.NewNumericAttribute("name", "place3"),
				*dynamodb.NewNumericAttribute("description", "place3 description"),
				*dynamodb.NewNumericAttribute("location", "41.558915,-8.131091"),
			},
			[]dynamodb.Attribute{
				*dynamodb.NewNumericAttribute("name", "place4"),
				*dynamodb.NewNumericAttribute("description", "place4 description"),
				*dynamodb.NewNumericAttribute("location", "42.458915,-6.631091"),
			},
		},
	},
	HashKeys: map[string][]string{
		"game_scores_test": []string{
			"Game X",
			"Game Y",
			"Game Y",
			"Game Z",
		},
	},
}

var _ = Suite(dynamoBaseRepository_suite)

func (s *DynamoBaseRepositorySuite) TestGetAllFromTableWithData(c *C) {
	// Arrange
	var repo DynamoBaseRepository

	// Act
	obtained, _ := repo.GetAll("users_test")

	// Assert
	expected := []map[string]string{
		{
			"id":         "3",
			"first_name": "Nuno",
			"last_name":  "Correia",
			"email":      "nuno@exemple.com",
			"country":    "Portugal",
		},
		{
			"id":         "2",
			"first_name": "Ana",
			"last_name":  "Dias",
			"email":      "ana@gmail.com",
			"country":    "Portugal",
		},
		{
			"id":         "1",
			"first_name": "Monika",
			"last_name":  "Sobkowiak",
			"email":      "monika@gmail.com",
			"country":    "Poland",
		},
	}

	c.Check(obtained, DeepEquals, expected)
}

func (s *DynamoBaseRepositorySuite) TestGetAllFromEmptyTable(c *C) {
	// Arrange
	var repo DynamoBaseRepository

	// Act
	obtained, _ := repo.GetAll("test")

	// Assert
	c.Check(len(obtained), Equals, 0)
}

func (s *DynamoBaseRepositorySuite) TestGetAllFromTableNonExistedTable(c *C) {
	// Arrange
	var repo DynamoBaseRepository

	// Arrange
	_, err := repo.GetAll("not_existed_table")

	// Assert
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")
}

func (s *DynamoBaseRepositorySuite) TestGetByHash(c *C) {
	var repo DynamoBaseRepository
	expected := map[string]string{
		"id":         "2",
		"first_name": "Ana",
		"last_name":  "Dias",
		"email":      "ana@gmail.com",
		"country":    "Portugal",
	}
	obtained, _ := repo.GetByHash("users_test", "2")

	c.Check(obtained, DeepEquals, expected)
}

func (s *DynamoBaseRepositorySuite) TestGetByHashRange(c *C) {
	var repo DynamoBaseRepository
	expected := map[string]string{
		"user_id":    "2",
		"wins":       "8",
		"game_title": "Game Y",
		"losts":      "0",
		"top_score":  "123",
	}
	obtained, _ := repo.GetByHashRange("game_scores_test", "Game Y", "2")

	c.Check(obtained, DeepEquals, expected)

	_, err := repo.GetByHashRange("not_existed_table", "1", "2")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	_, err = repo.GetByHashRange("game_scores_test", "Game Y", "10")
	c.Check(err, ErrorMatches, "Item not found")
}

func (s *DynamoBaseRepositorySuite) TestGetByOnlyHash(c *C) {
	var repo DynamoBaseRepository
	expected := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtained, _ := repo.GetByOnlyHash("game_scores_test", "Game Y")
	c.Check(obtained, DeepEquals, expected)

	_, err := repo.GetByOnlyHash("not_existed_table", "1")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	obtained, _ = repo.GetByOnlyHash("game_scores_test", "not_existed_hash")
	expected = []map[string]string{}
	c.Check(obtained, DeepEquals, expected)
}

func (s *DynamoBaseRepositorySuite) TestGetByIndexHash(c *C) {
	var repo DynamoBaseRepository
	expected := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "582",
			"user_id":    "1",
			"wins":       "8",
			"game_title": "Game X",
			"losts":      "2",
		},
		map[string]string{
			"top_score":  "12",
			"user_id":    "4",
			"wins":       "8",
			"game_title": "Game Z",
			"losts":      "2",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtained, _ := repo.GetByIndexHash("game_scores_test", "wins_losts", "8")
	c.Check(obtained, DeepEquals, expected)

	_, err := repo.GetByIndexHash("not_existed_table", "index", "1")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	obtained, _ = repo.GetByIndexHash("game_scores_test", "not_existed_index", "1")
	expected = []map[string]string(nil)
	c.Check(obtained, DeepEquals, expected)

	expected = []map[string]string{
		map[string]string{
			"country":    "Poland",
			"email":      "monika@gmail.com",
			"first_name": "Monika",
			"id":         "1",
			"last_name":  "Sobkowiak",
		},
	}
	obtained, _ = repo.GetByIndexHash("users_test", "email", "monika@gmail.com")
	c.Check(obtained, DeepEquals, expected)
}

func (s *DynamoBaseRepositorySuite) TestGetByOnlyRange(c *C) {
	var repo DynamoBaseRepository
	expectedGT := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedGT, _ := repo.GetByOnlyRange("game_scores_test", "Game Y", "GT", []string{"1"})
	c.Check(obtainedGT, DeepEquals, expectedGT)

	expectedGE := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedGE, _ := repo.GetByOnlyRange("game_scores_test", "Game Y", "GE", []string{"2"})
	c.Check(obtainedGE, DeepEquals, expectedGE)

	expectedLT := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
	}

	obtainedLT, _ := repo.GetByOnlyRange("game_scores_test", "Game Y", "LT", []string{"3"})
	c.Check(obtainedLT, DeepEquals, expectedLT)

	expectedLE := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedLE, _ := repo.GetByOnlyRange("game_scores_test", "Game Y", "LE", []string{"3"})
	c.Check(obtainedLE, DeepEquals, expectedLE)

	expectedBETWEEN := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedBETWEEN, _ := repo.GetByOnlyRange("game_scores_test", "Game Y", "BETWEEN", []string{"2", "4"})
	c.Check(obtainedBETWEEN, DeepEquals, expectedBETWEEN)
}

func (s *DynamoBaseRepositorySuite) TestGetByIndexRange(c *C) {
	var repo DynamoBaseRepository
	expectedGT := []map[string]string{
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedGT, _ := repo.GetByIndexRange("game_scores_test", "wins_losts", "8", "GT", []string{"2"})
	c.Check(obtainedGT, DeepEquals, expectedGT)

	expectedGE := []map[string]string{
		map[string]string{
			"top_score":  "582",
			"user_id":    "1",
			"wins":       "8",
			"game_title": "Game X",
			"losts":      "2",
		},
		map[string]string{
			"top_score":  "12",
			"user_id":    "4",
			"wins":       "8",
			"game_title": "Game Z",
			"losts":      "2",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedGE, _ := repo.GetByIndexRange("game_scores_test", "wins_losts", "8", "GE", []string{"2"})
	c.Check(obtainedGE, DeepEquals, expectedGE)

	expectedLT := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
	}

	obtainedLT, _ := repo.GetByIndexRange("game_scores_test", "wins_losts", "8", "LT", []string{"2"})
	c.Check(obtainedLT, DeepEquals, expectedLT)

	expectedLE := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "582",
			"user_id":    "1",
			"wins":       "8",
			"game_title": "Game X",
			"losts":      "2",
		},
		map[string]string{
			"top_score":  "12",
			"user_id":    "4",
			"wins":       "8",
			"game_title": "Game Z",
			"losts":      "2",
		},
	}

	obtainedLE, _ := repo.GetByIndexRange("game_scores_test", "wins_losts", "8", "LE", []string{"2"})
	c.Check(obtainedLE, DeepEquals, expectedLE)

	expectedBETWEEN := []map[string]string{
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedBETWEEN, _ := repo.GetByIndexRange("game_scores_test", "wins_losts", "8", "BETWEEN", []string{"80", "100"})
	c.Check(obtainedBETWEEN, DeepEquals, expectedBETWEEN)

}

func (s *DynamoBaseRepositorySuite) TestAdd(c *C) {
	var repo DynamoBaseRepository
	item := []Attribute{
		Attribute{
			Description: AttributeDefinition{Name: "first_name", Type: "S"},
			Value:       "Test_first_name",
		},
		Attribute{
			Description: AttributeDefinition{Name: "last_name", Type: "S"},
			Value:       "Test_last_name",
		},
		Attribute{
			Description: AttributeDefinition{Name: "email", Type: "S"},
			Value:       "Test_email",
		},
		Attribute{
			Description: AttributeDefinition{Name: "country", Type: "S"},
			Value:       "Test_country",
		},
	}

	status, _ := repo.Add("users_test", "4", "", item)
	c.Check(status, Equals, true)

	expected := map[string]string{
		"country":    "Test_country",
		"email":      "Test_email",
		"first_name": "Test_first_name",
		"id":         "4",
		"last_name":  "Test_last_name",
	}

	obtained, _ := repo.GetByHash("users_test", "4")
	c.Check(obtained, DeepEquals, expected)
}

func (s *DynamoBaseRepositorySuite) TestDeleteByHash(c *C) {
	var repo DynamoBaseRepository
	status, _ := repo.DeleteByHash("users_test", "4")
	c.Check(status, Equals, true)

	_, err := repo.GetByHash("users_test", "4")
	c.Check(err, ErrorMatches, "Item not found")
}

func (s *DynamoBaseRepositorySuite) TestAddByHashRange(c *C) {
	var repo DynamoBaseRepository
	item := []Attribute{
		Attribute{
			Description: AttributeDefinition{Name: "top_score", Type: "N"},
			Value:       "111",
		},
		Attribute{
			Description: AttributeDefinition{Name: "wins", Type: "N"},
			Value:       "222",
		},
		Attribute{
			Description: AttributeDefinition{Name: "losts", Type: "N"},
			Value:       "333",
		},
	}

	status, _ := repo.Add("game_scores_test", "test_hash_value", "1", item)
	c.Check(status, Equals, true)

	expected := map[string]string{
		"top_score":  "111",
		"user_id":    "1",
		"wins":       "222",
		"game_title": "test_hash_value",
		"losts":      "333",
	}

	obtained, _ := repo.GetByHashRange("game_scores_test", "test_hash_value", "1")
	c.Check(obtained, DeepEquals, expected)
}

func (s *DynamoBaseRepositorySuite) TestDeleteByRange(c *C) {
	var repo DynamoBaseRepository
	status, _ := repo.DeleteByHashRange("game_scores_test", "test_hash_value", "1")
	c.Check(status, Equals, true)

	_, err := repo.GetByHashRange("game_scores_test", "test_range_value", "1")
	c.Check(err, ErrorMatches, "Item not found")
}
