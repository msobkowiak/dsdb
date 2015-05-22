package main

import (
	//"github.com/goamz/goamz/dynamodb"
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

	obtained, _ := RepoGetAllItems("users")

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
	obtained, _ := RepoGetItemByHash("users", "2")

	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestRepoGetItemByHashRange(c *C) {
	expected := map[string]string{
		"user_id":    "2",
		"wins":       "3",
		"game_title": "Game Y",
		"losts":      "0",
		"top_score":  "123",
	}
	obtained, _ := RepoGetItemByHashRange("game_scores", "Game Y", "2")

	c.Check(obtained, DeepEquals, expected)

	_, err := RepoGetItemByHashRange("not_existed_table", "1", "2")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	_, err = RepoGetItemByHashRange("game_scores", "Game Y", "10")
	c.Check(err, ErrorMatches, "Item not found")
}

func (s *TableSuite) TestRepoGetItemsByHash(c *C) {
	expected := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "3",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "30",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtained, _ := RepoGetItemsByHash("game_scores", "Game Y")
	c.Check(obtained, DeepEquals, expected)

	_, err := RepoGetItemsByHash("not_existed_table", "1")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	obtained, _ = RepoGetItemsByHash("game_scores", "not_existed_hash")
	expected = []map[string]string{}
	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestRepoGetItemByIndexHash(c *C) {
	expected := []map[string]string{
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "30",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtained, _ := RepoGetItemByIndexHash("game_scores", "wins_losts", "30")
	c.Check(obtained, DeepEquals, expected)

	_, err := RepoGetItemByIndexHash("not_existed_table", "index", "1")
	c.Check(err, ErrorMatches, "Table not_existed_table not found.")

	obtained, _ = RepoGetItemByIndexHash("game_scores", "not_existed_index", "1")
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
	obtained, _ = RepoGetItemByIndexHash("users", "email", "monika@gmail.com")
	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestRepoGetItemsByRangeOp(c *C) {
	expectedGT := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "3",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "30",
			"game_title": "Game Y",
			"losts":      "90",
		},
		/*map[string]string{
			"top_score":  "12",
			"user_id":    "4",
			"wins":       "2",
			"game_title": "Game Z",
			"losts":      "2",
		},*/
	}

	obtainedGT, _ := RepoGetItemsByRangeOp("game_scores", "Game Y", "GT", []string{"1"})
	c.Check(obtainedGT, DeepEquals, expectedGT)

	expectedGE := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "3",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "30",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedGE, _ := RepoGetItemsByRangeOp("game_scores", "Game Y", "GE", []string{"2"})
	c.Check(obtainedGE, DeepEquals, expectedGE)

	expectedLT := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "3",
			"game_title": "Game Y",
			"losts":      "0",
		},
	}

	obtainedLT, _ := RepoGetItemsByRangeOp("game_scores", "Game Y", "LT", []string{"3"})
	c.Check(obtainedLT, DeepEquals, expectedLT)

	expectedLE := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "3",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "30",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedLE, _ := RepoGetItemsByRangeOp("game_scores", "Game Y", "LE", []string{"3"})
	c.Check(obtainedLE, DeepEquals, expectedLE)

	expectedBETWEEN := []map[string]string{
		map[string]string{
			"top_score":  "123",
			"user_id":    "2",
			"wins":       "3",
			"game_title": "Game Y",
			"losts":      "0",
		},
		map[string]string{
			"top_score":  "333333",
			"user_id":    "3",
			"wins":       "30",
			"game_title": "Game Y",
			"losts":      "90",
		},
	}

	obtainedBETWEEN, _ := RepoGetItemsByRangeOp("game_scores", "Game Y", "BETWEEN", []string{"2", "4"})
	c.Check(obtainedBETWEEN, DeepEquals, expectedBETWEEN)

}

//func RepoGetItemsByRangeOp(tableName, hashValue, operator string, rangeValue []string) ([]map[string]string, error)
