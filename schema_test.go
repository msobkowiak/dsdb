package main

import (
	. "gopkg.in/check.v1"
)

func (s *TableSuite) TestHasRange(c *C) {
	c.Check(table_suite.Db.Tables["users"].HasRange(), Equals, false)
	c.Check(table_suite.Db.Tables["game_scores"].HasRange(), Equals, true)
}

func (s *TableSuite) TestGetTypeOfAttribute(c *C) {
	c.Check(table_suite.Db.Tables["users"].GetTypeOfAttribute("id"), Equals, "N")
	c.Check(table_suite.Db.Tables["users"].GetTypeOfAttribute("email"), Equals, "S")
}

func (s *TableSuite) TestGetIndexByName(c *C) {
	expexted := SecondaryIndexDefinition{
		Name: "country",
		Type: "HASH",
		Hash: "country",
	}

	obtained, _ := table_suite.Db.Tables["users"].GetIndexByName("country")
	_, err := table_suite.Db.Tables["users"].GetIndexByName("not_existed_index")
	c.Check(obtained, DeepEquals, expexted)
	c.Check(err, ErrorMatches, "Index not found")
}
