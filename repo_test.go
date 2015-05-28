package main

import (
	. "gopkg.in/check.v1"
)

func (s *TableSuite) TestRepoGetAllItems(c *C) {
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

	obtained, _ := RepoGetAllItems("users_test")

	c.Check(obtained, DeepEquals, expected)

	_, err := RepoGetAllItems("not_existed_table")

	c.Check(err, ErrorMatches, "Table not_existed_table not found.")
}

func (s *TableSuite) TestRepoGetItemByHash(c *C) {
	expected := map[string]string{
		"id":         "2",
		"first_name": "Ana",
		"last_name":  "Dias",
		"email":      "ana@gmail.com",
		"country":    "Portugal",
	}
	obtained, _ := RepoGetItemByHash("users_test", "2")

	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestRepoGetItemByHashRange(c *C) {
	expected := map[string]string{
		"user_id":    "2",
		"wins":       "8",
		"game_title": "Game Y",
		"losts":      "0",
		"top_score":  "123",
	}
	obtained, _ := RepoGetItemByHashRange("game_scores_test", "Game Y", "2")

	c.Check(obtained, DeepEquals, expected)

	_, err := RepoGetItemByHashRange("not_existed_table", "1", "2")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	_, err = RepoGetItemByHashRange("game_scores_test", "Game Y", "10")
	c.Check(err, ErrorMatches, "Item not found")
}

func (s *TableSuite) TestRepoGetItemsByHash(c *C) {
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

	obtained, _ := RepoGetItemsByHash("game_scores_test", "Game Y")
	c.Check(obtained, DeepEquals, expected)

	_, err := RepoGetItemsByHash("not_existed_table", "1")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	obtained, _ = RepoGetItemsByHash("game_scores_test", "not_existed_hash")
	expected = []map[string]string{}
	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestRepoGetItemByIndexHash(c *C) {
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

	obtained, _ := RepoGetItemByIndexHash("game_scores_test", "wins_losts", "8")
	c.Check(obtained, DeepEquals, expected)

	_, err := RepoGetItemByIndexHash("not_existed_table", "index", "1")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	obtained, _ = RepoGetItemByIndexHash("game_scores_test", "not_existed_index", "1")
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
	obtained, _ = RepoGetItemByIndexHash("users_test", "email", "monika@gmail.com")
	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestRepoGetItemsByRangeOp(c *C) {
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

	obtainedGT, _ := RepoGetItemsByRangeOp("game_scores_test", "Game Y", "GT", []string{"1"})
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

	obtainedGE, _ := RepoGetItemsByRangeOp("game_scores_test", "Game Y", "GE", []string{"2"})
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

	obtainedLT, _ := RepoGetItemsByRangeOp("game_scores_test", "Game Y", "LT", []string{"3"})
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

	obtainedLE, _ := RepoGetItemsByRangeOp("game_scores_test", "Game Y", "LE", []string{"3"})
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

	obtainedBETWEEN, _ := RepoGetItemsByRangeOp("game_scores_test", "Game Y", "BETWEEN", []string{"2", "4"})
	c.Check(obtainedBETWEEN, DeepEquals, expectedBETWEEN)
}

func (s *TableSuite) TestRepoGetItemsByIndexRangeOp(c *C) {
	expectedGT := []map[string]string{
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "8",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedGT, _ := RepoGetItemsByIndexRangeOp("game_scores_test", "wins_losts", "8", "GT", []string{"2"})
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

	obtainedGE, _ := RepoGetItemsByIndexRangeOp("game_scores_test", "wins_losts", "8", "GE", []string{"2"})
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

	obtainedLT, _ := RepoGetItemsByIndexRangeOp("game_scores_test", "wins_losts", "8", "LT", []string{"2"})
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

	obtainedLE, _ := RepoGetItemsByIndexRangeOp("game_scores_test", "wins_losts", "8", "LE", []string{"2"})
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

	obtainedBETWEEN, _ := RepoGetItemsByIndexRangeOp("game_scores_test", "wins_losts", "8", "BETWEEN", []string{"80", "100"})
	c.Check(obtainedBETWEEN, DeepEquals, expectedBETWEEN)

}

func (s *TableSuite) TestAddItem(c *C) {
	item := []Attribute{
		Attribute{
			Description: AttributeDefinition{
				Name: "first_name",
				Type: "S",
			},
			Value: "Test_first_name",
		},
		Attribute{
			Description: AttributeDefinition{
				Name: "last_name",
				Type: "S",
			},
			Value: "Test_last_name",
		},
		Attribute{
			Description: AttributeDefinition{
				Name: "email",
				Type: "S",
			},
			Value: "Test_email",
		},
		Attribute{
			Description: AttributeDefinition{
				Name: "country",
				Type: "S",
			},
			Value: "Test_country",
		},
	}

	status, _ := RepoAddItem("users_test", "4", item)
	c.Check(status, Equals, true)

	expected := map[string]string{
		"country":    "Test_country",
		"email":      "Test_email",
		"first_name": "Test_first_name",
		"id":         "4",
		"last_name":  "Test_last_name",
	}

	obtained, _ := RepoGetItemByHash("users_test", "4")
	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestRepoDeleteItem(c *C) {
	status, _ := RepoDeleteItem("users_test", "4")
	c.Check(status, Equals, true)

	_, err := RepoGetItemByHash("users_test", "4")
	c.Check(err, ErrorMatches, "Item not found")
}

func (s *TableSuite) TestAddItemHashRange(c *C) {
	item := []Attribute{
		Attribute{
			Description: AttributeDefinition{
				Name: "top_score",
				Type: "N",
			},
			Value: "111",
		},
		Attribute{
			Description: AttributeDefinition{
				Name: "wins",
				Type: "N",
			},
			Value: "222",
		},
		Attribute{
			Description: AttributeDefinition{
				Name: "losts",
				Type: "N",
			},
			Value: "333",
		},
	}

	status, _ := RepoAddItemHashRange("game_scores_test", "test_hash_value", "1", item)
	c.Check(status, Equals, true)

	expected := map[string]string{
		"top_score":  "111",
		"user_id":    "1",
		"wins":       "222",
		"game_title": "test_hash_value",
		"losts":      "333",
	}

	obtained, _ := RepoGetItemByHashRange("game_scores_test", "test_hash_value", "1")
	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestRepoDeleteItemWithRange(c *C) {
	status, _ := RepoDeleteItemWithRange("game_scores_test", "test_hash_value", "1")
	c.Check(status, Equals, true)

	_, err := RepoGetItemByHashRange("game_scores_test", "test_range_value", "1")
	c.Check(err, ErrorMatches, "Item not found")
}
