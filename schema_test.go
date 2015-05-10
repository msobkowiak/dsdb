package main

import (
	. "gopkg.in/check.v1"
)

func (s *TableSuite) TestHasRange(c *C) {
	c.Check(table_suite.Tables["users"].HasRange(), Equals, false)
	c.Check(table_suite.Tables["game_scores"].HasRange(), Equals, true)
}

func (s *TableSuite) TestGetTypeOfAttribute(c *C) {
	c.Check(table_suite.Tables["users"].GetTypeOfAttribute("id"), Equals, "N")
	c.Check(table_suite.Tables["users"].GetTypeOfAttribute("email"), Equals, "S")
}
