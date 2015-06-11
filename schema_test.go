package main

import (
	. "gopkg.in/check.v1"
)

func (s *TableSuite) TestHasRange(c *C) {
	c.Check(table_suite.Db.Tables["users_test"].HasRange(), Equals, false)
	c.Check(table_suite.Db.Tables["game_scores_test"].HasRange(), Equals, true)
}

func (s *TableSuite) TestGetTypeOfAttribute(c *C) {
	c.Check(table_suite.Db.Tables["users_test"].GetTypeOfAttribute("id"), Equals, "N")
	c.Check(table_suite.Db.Tables["users_test"].GetTypeOfAttribute("email"), Equals, "S")
	c.Check(table_suite.Db.Tables["game_scores_test"].GetTypeOfAttribute("wins"), Equals, "N")
	c.Check(table_suite.Db.Tables["users_test"].GetTypeOfAttribute("not_existed_atrr"), Equals, "S")
}

func (s *TableSuite) TestGetIndexByName(c *C) {
	expexted := SecondaryIndexDefinition{
		Name: "country",
		Type: "HASH",
		Hash: "country",
	}

	obtained, _ := table_suite.Db.Tables["users_test"].GetIndexByName("country")
	_, err := table_suite.Db.Tables["users"].GetIndexByName("not_existed_index")
	c.Check(obtained, DeepEquals, expexted)
	c.Check(err, ErrorMatches, "Index not found")
}

func (s *TableSuite) TestGetHashName(c *C) {
	expected := "id"
	obtained, _ := GetHashName("users_test", table_suite.Db)
	c.Check(obtained, Equals, expected)

	_, err := GetHashName("not_existet_table", table_suite.Db)
	c.Check(err, ErrorMatches, "Table not_existet_table not found.")
}

func (s *TableSuite) TestGetRangeName(c *C) {
	expected := "user_id"
	obtained, _ := GetRangeName("game_scores_test", table_suite.Db)
	c.Check(obtained, Equals, expected)

	obtained, _ = GetRangeName("users_test", table_suite.Db)
	c.Check(obtained, Equals, "")

	_, err := GetHashName("not_existet_table", table_suite.Db)
	c.Check(err, ErrorMatches, "Table not_existet_table not found.")
}
