package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TableSuite struct {
	Tables   map[string]TableDescription
	Data     map[string][][]dynamodb.Attribute
	HashKeys map[string][]string
}

func (s *TableSuite) SetUpSuite(c *C) {
	users := CreateTable(table_suite.Tables["users"])
	games := CreateTable(table_suite.Tables["game_scores"])
	AddItems(users, table_suite.Data["users"], nil)
	AddItems(games, table_suite.Data["game_scores"], table_suite.HashKeys["game_scores"])
}

func (s *TableSuite) TearDownSuite(c *C) {
	dynamAuth := table_suite.Tables["users"].Authentication.Dynamo
	db := Auth(dynamAuth.Region, dynamAuth.AccessKey, dynamAuth.SecretKey)
	DeleteAllTables(db)
}

var table_suite = &TableSuite{
	Tables: map[string]TableDescription{
		"users": TableDescription{
			Name: "users",
			Attributes: []AttributeDefinition{
				AttributeDefinition{"id", "N", true},
				AttributeDefinition{"email", "S", true},
				AttributeDefinition{"country", "S", true},
			},
			PrimaryKey: PrimaryKeyDefinition{
				Type: "HASH",
				Hash: "id",
			},
			SecondaryIndexes: []SecondaryIndexDefinition{
				SecondaryIndexDefinition{
					Name: "email",
					Type: "HASH",
					Hash: "email",
				},
				SecondaryIndexDefinition{
					Name: "country",
					Type: "HASH",
					Hash: "country",
				},
			},
			Authentication: Authentication{
				DynamoAuth{
					Region:    "http://127.0.0.1:4567",
					AccessKey: "access",
					SecretKey: "secret",
				},
			},
		},
		"game_scores": TableDescription{
			Name: "game_scores",
			Attributes: []AttributeDefinition{
				AttributeDefinition{"user_id", "N", true},
				AttributeDefinition{"game_title", "S", true},
				AttributeDefinition{"wins", "N", true},
				AttributeDefinition{"losts", "N", true},
			},
			PrimaryKey: PrimaryKeyDefinition{
				Type:  "RANGE",
				Hash:  "game_title",
				Range: "user_id",
			},
			SecondaryIndexes: []SecondaryIndexDefinition{
				SecondaryIndexDefinition{
					Name:  "wins_losts",
					Type:  "RANGE",
					Hash:  "wins",
					Range: "losts",
				},
			},
			Authentication: Authentication{
				DynamoAuth{
					Region:    "http://127.0.0.1:4567",
					AccessKey: "access",
					SecretKey: "secret",
				},
			},
		},
	},
	Data: map[string][][]dynamodb.Attribute{
		"users": [][]dynamodb.Attribute{
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
		"game_scores": [][]dynamodb.Attribute{
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
	},
	HashKeys: map[string][]string{
		"game_scores": []string{
			"Game X",
			"Game Y",
			"Game Y",
			"Game Z",
		},
	},
}

var _ = Suite(table_suite)

func (s *TableSuite) TestAuth(c *C) {
	region := "http://127.0.0.1:4567"
	accessKey := "key"
	secretKey := "secret"
	obtained := Auth(region, accessKey, secretKey)

	expected := dynamodb.Server{
		Auth:   aws.Auth{AccessKey: "key", SecretKey: "secret"},
		Region: aws.Region{DynamoDBEndpoint: "http://127.0.0.1:4567"},
	}

	c.Check(obtained, Equals, expected)
}

func (s *TableSuite) TestConvertToDynamo(c *C) {
	var expected = dynamodb.TableDescriptionT{
		TableName: "users",
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

	obtained := ConvertToDynamo(table_suite.Tables["users"])

	c.Check(obtained, DeepEquals, expected)
}

func (s *TableSuite) TestGetDynamoTable(c *C) {
	expected := dynamodb.Table{
		Server: &dynamodb.Server{
			Auth:   aws.Auth{AccessKey: "access", SecretKey: "secret"},
			Region: aws.Region{DynamoDBEndpoint: "http://127.0.0.1:4567"},
		},
		Name: "users",
		Key: dynamodb.PrimaryKey{
			KeyAttribute: &dynamodb.Attribute{
				Type: "N",
				Name: "id",
			},
		},
	}
	obtained, _ := GetDynamoTable("users")

	c.Check(obtained, DeepEquals, expected)

	_, err := GetDynamoTable("not_existed_table")

	c.Check(err, ErrorMatches, "Table not_existed_table not found.")
}
