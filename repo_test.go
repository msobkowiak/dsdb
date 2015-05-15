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

	//_, err := RepoGetAllItems("not_existed_table")

	//c.Check(err, ErrorMatches, "Table not_existed_table does not exists.")
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
}
