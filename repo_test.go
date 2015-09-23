package main

import (
	. "gopkg.in/check.v1"
)

func (s *TableSuite) TestGetAll(c *C) {
	var repo DynamoBaseRepository
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

	obtained, _ := repo.GetAll("users_test")

	c.Check(obtained, DeepEquals, expected)

	_, err := repo.GetAll("not_existed_table")

	c.Check(err, ErrorMatches, "Table not_existed_table not found.")
}

func (s *TableSuite) TestGetByHash(c *C) {
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

func (s *TableSuite) TestGetByHashRange(c *C) {
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

func (s *TableSuite) TestGetByOnlyHash(c *C) {
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

func (s *TableSuite) TestGetByIndexHash(c *C) {
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

func (s *TableSuite) TestGetByOnlyRange(c *C) {
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

func (s *TableSuite) TestGetByIndexRange(c *C) {
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

func (s *TableSuite) TestAdd(c *C) {
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

func (s *TableSuite) TestDeleteByHash(c *C) {
	var repo DynamoBaseRepository
	status, _ := repo.DeleteByHash("users_test", "4")
	c.Check(status, Equals, true)

	_, err := repo.GetByHash("users_test", "4")
	c.Check(err, ErrorMatches, "Item not found")
}

func (s *TableSuite) TestAddByHashRange(c *C) {
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

func (s *TableSuite) TestDeleteByRange(c *C) {
	var repo DynamoBaseRepository
	status, _ := repo.DeleteByHashRange("game_scores_test", "test_hash_value", "1")
	c.Check(status, Equals, true)

	_, err := repo.GetByHashRange("game_scores_test", "test_range_value", "1")
	c.Check(err, ErrorMatches, "Item not found")
}
