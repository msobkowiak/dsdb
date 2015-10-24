package main

import (
	. "gopkg.in/check.v1"
)

type TableDescriptionSuite struct {
	Db DbDescription
}

var tables = tableDescription_suite.Db.Tables
var dbDescription = tableDescription_suite.Db

var tableDescription_suite = &TableDescriptionSuite{

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
			"geo_test": TableDescription{
				Name: "geo_test",
				Attributes: []AttributeDefinition{
					AttributeDefinition{"id", "N", true},
					AttributeDefinition{"name", "S", true},
					AttributeDefinition{"descroption", "S", true},
					AttributeDefinition{"location", "G", true},
				},
				PrimaryKey: KeyDefinition{
					Type: "HASH",
					Hash: "id",
				},
			},
		},
	},
}

var _ = Suite(tableDescription_suite)

func (s *TableDescriptionSuite) TestHasRange(c *C) {
	c.Check(tables["users_test"].HasRange(), Equals, false)
	c.Check(tables["game_scores_test"].HasRange(), Equals, true)
}

func (s *TableDescriptionSuite) TestHasGeoPoint(c *C) {
	c.Check(tables["geo_test"].HasGeoPoint(), Equals, true)
	c.Check(tables["game_scores_test"].HasGeoPoint(), Equals, false)
}

func (s *TableDescriptionSuite) TestGetTypeOfAttribute(c *C) {
	c.Check(tables["users_test"].GetTypeOfAttribute("id"), Equals, "N")
	c.Check(tables["users_test"].GetTypeOfAttribute("email"), Equals, "S")
	c.Check(tables["game_scores_test"].GetTypeOfAttribute("wins"), Equals, "N")
	c.Check(tables["users_test"].GetTypeOfAttribute("not_existed_atrr"), Equals, "S")
}

func (s *TableDescriptionSuite) TestGetIndexByNameReturrnsAnIndex(c *C) {
	// Arrange
	var expectedSecondaryIndex = SecondaryIndexDefinition{
		Name: "country",
		Key: KeyDefinition{
			Type: "HASH",
			Hash: "country",
		},
	}
	var othreSecondaryIndex = SecondaryIndexDefinition{
		Name: "test",
		Key: KeyDefinition{
			Type: "HASH",
			Hash: "test",
		},
	}
	var table = TableDescription{
		Name: "users_test",
		SecondaryIndexes: []SecondaryIndexDefinition{
			othreSecondaryIndex,
			expectedSecondaryIndex,
		},
	}

	// Act
	obtained, _ := table.GetIndexByName("country")

	// Assert
	c.Check(obtained, DeepEquals, expectedSecondaryIndex)
}

func (s *TableDescriptionSuite) TestGetIndexByNameReturrnsAnError(c *C) {
	// Arrange
	var expectedSecondaryIndex = SecondaryIndexDefinition{
		Name: "country",
		Key: KeyDefinition{
			Type: "HASH",
			Hash: "country",
		},
	}
	var othreSecondaryIndex = SecondaryIndexDefinition{
		Name: "test",
		Key: KeyDefinition{
			Type: "HASH",
			Hash: "test",
		},
	}
	var table = TableDescription{
		Name: "users_test",
		SecondaryIndexes: []SecondaryIndexDefinition{
			othreSecondaryIndex,
			expectedSecondaryIndex,
		},
	}

	// Act
	_, err := table.GetIndexByName("non_existed_index_name")

	// Assert
	c.Check(err, ErrorMatches, "Index not found")
}

func (s *TableDescriptionSuite) TestGeGeoPointName(c *C) {
	expected := "location"
	obtained, _ := tables["geo_test"].GetGeoPointName()
	c.Check(obtained, Equals, expected)

	_, err := tables["users_test"].GetGeoPointName()
	c.Check(err, ErrorMatches, "No geo point found")
}

func (s *TableDescriptionSuite) TestExcludeNonKeyAttributes(c *C) {
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
		PrimaryKey: KeyDefinition{
			Hash:  "attr4",
			Range: "attr1",
		},
		SecondaryIndexes: []SecondaryIndexDefinition{
			SecondaryIndexDefinition{Key: KeyDefinition{Hash: "attr2", Range: "attr3"}},
			SecondaryIndexDefinition{Key: KeyDefinition{Hash: "attr3"}},
		},
	}

	obtained := testData.ExcludeNonKeyAttributes()
	c.Check(obtained, DeepEquals, expected)
}

func (s *TableDescriptionSuite) TestGetRangeName(c *C) {
	expected := "user_id"
	obtained, _ := dbDescription.GetRangeName("game_scores_test")
	c.Check(obtained, Equals, expected)

	obtained, _ = dbDescription.GetRangeName("users_test")
	c.Check(obtained, Equals, "")

	_, err := dbDescription.GetHashName("not_existet_table")
	c.Check(err, ErrorMatches, "Table not_existet_table not found.")
}

func (s *TableDescriptionSuite) TestGetHashName(c *C) {
	expected := "id"
	obtained, _ := dbDescription.GetHashName("users_test")
	c.Check(obtained, Equals, expected)

	_, err := dbDescription.GetHashName("not_existet_table")
	c.Check(err, ErrorMatches, "Table not_existet_table not found.")
}
