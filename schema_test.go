package main

import (
	. "gopkg.in/check.v1"
)

func (s *TableSuite) TestHasRange(c *C) {
	c.Check(table_suite.Db.Tables["users_test"].HasRange(), Equals, false)
	c.Check(table_suite.Db.Tables["game_scores_test"].HasRange(), Equals, true)
}

func (s *TableSuite) TestHasGeoPoint(c *C) {
	c.Check(table_suite.Db.Tables["geo_test"].HasGeoPoint(), Equals, true)
	c.Check(table_suite.Db.Tables["game_scores_test"].HasGeoPoint(), Equals, false)
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
	_, err := table_suite.Db.Tables["users_test"].GetIndexByName("not_existed_index")
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

func (s *TableSuite) TestGeGeoPointName(c *C) {
	expected := "location"
	obtained, _ := table_suite.Db.Tables["geo_test"].GetGeoPointName()
	c.Check(obtained, Equals, expected)

	_, err := table_suite.Db.Tables["users_test"].GetGeoPointName()
	c.Check(err, ErrorMatches, "No geo point found")
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

func (s *TableSuite) TestExcludeNonKeyAttributes(c *C) {
	expected := []AttributeDefinition{
		AttributeDefinition{Name: "attr1", Type: "S", Required: false},
		AttributeDefinition{Name: "attr2", Type: "N", Required: true},
		AttributeDefinition{Name: "attr3", Type: "S", Required: true},
		AttributeDefinition{Name: "attr4", Type: "S", Required: true},
	}

	testData := TableDescription{
		Attributes: []AttributeDefinition{
			AttributeDefinition{Name: "id", Type: "N", Required: true},
			AttributeDefinition{Name: "attr1", Type: "S", Required: false},
			AttributeDefinition{Name: "attr2", Type: "N", Required: true},
			AttributeDefinition{Name: "attr3", Type: "S", Required: true},
			AttributeDefinition{Name: "attr4", Type: "S", Required: true},
		},
		PrimaryKey: PrimaryKeyDefinition{
			Hash:  "attr4",
			Range: "attr1",
		},
		SecondaryIndexes: []SecondaryIndexDefinition{
			SecondaryIndexDefinition{Hash: "attr2", Range: "attr3"},
			SecondaryIndexDefinition{Hash: "attr3"},
		},
	}

	obtained := ExcludeNonKeyAttributes(testData)
	c.Check(obtained, DeepEquals, expected)
}
