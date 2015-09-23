package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var table DynamoTable

type TableSuite struct {
	Db       DbDescription
	Data     map[string][][]dynamodb.Attribute
	HashKeys map[string][]string
}

func (s *TableSuite) SetUpSuite(c *C) {
	schema = table_suite.Db
	var t DynamoTable

	users := t.Create(table_suite.Db.Tables["users_test"])
	games := t.Create(table_suite.Db.Tables["game_scores_test"])
	t.AddItems(users, table_suite.Data["users_test"], nil)
	t.AddItems(games, table_suite.Data["game_scores_test"], table_suite.HashKeys["game_scores_test"])
}

func (s *TableSuite) TearDownSuite(c *C) {
	var client DynamoClient
	dynamAuth := table_suite.Db.Authentication.Dynamo
	db := client.Auth(dynamAuth.Region, dynamAuth.AccessKey, dynamAuth.SecretKey)
	client.DeleteAll(db)
}

var table_suite = &TableSuite{

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

var _ = Suite(table_suite)

func (s *TableSuite) MapTest(c *C) {
	var expected = dynamodb.TableDescriptionT{
		TableName: "users_test",
		AttributeDefinitions: []dynamodb.AttributeDefinitionT{
			dynamodb.AttributeDefinitionT{"id", "N"},
			dynamodb.AttributeDefinitionT{"email", "S"},
			dynamodb.AttributeDefinitionT{"country", "S"},
		},
		KeySchema: []dynamodb.KeySchemaT{
			dynamodb.KeySchemaT{"id", "HASH"},
		},
		GlobalSecondaryIndexes: []dynamodb.GlobalSecondaryIndexT{
			dynamodb.GlobalSecondaryIndexT{
				IndexName: "email",
				KeySchema: []dynamodb.KeySchemaT{
					dynamodb.KeySchemaT{"email", "HASH"},
				},
				Projection: dynamodb.ProjectionT{"ALL"},
				ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
					ReadCapacityUnits:  1,
					WriteCapacityUnits: 1,
				},
			},
			dynamodb.GlobalSecondaryIndexT{
				IndexName: "country",
				KeySchema: []dynamodb.KeySchemaT{
					dynamodb.KeySchemaT{"country", "HASH"},
				},
				Projection: dynamodb.ProjectionT{"ALL"},
				ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
					ReadCapacityUnits:  1,
					WriteCapacityUnits: 1,
				},
			},
		},
		ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
			ReadCapacityUnits:  10,
			WriteCapacityUnits: 10,
		},
	}

	obtained := table.Map(table_suite.Db.Tables["users_test"])

	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestGetByName(c *C) {
	expected := dynamodb.Table{
		Server: &dynamodb.Server{
			Auth:   aws.Auth{AccessKey: "access", SecretKey: "secret"},
			Region: aws.Region{DynamoDBEndpoint: "http://127.0.0.1:4567"},
		},
		Name: "users_test",
		Key: dynamodb.PrimaryKey{
			KeyAttribute: &dynamodb.Attribute{
				Type: "N",
				Name: "id",
			},
		},
	}
	obtained, _ := table.GetByName("users_test")

	c.Check(obtained, DeepEquals, expected)

	_, err := table.GetByName("not_existed_table")

	c.Check(err, ErrorMatches, "Table not_existed_table not found.")
}
