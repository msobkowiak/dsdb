package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TableSuite struct {
	Tables map[string]TableDescription
}

func (s *TableSuite) SetUpSuite(c *C) {
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
	obtained := GetDynamoTable("users")

	c.Check(obtained, DeepEquals, expected)
}
