package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var table DynamoTable

type DynamoTableSuite struct {
	Db DbDescription
}

var dynamoTable_suite = &DynamoTableSuite{

	Db: DbDescription{
		Name:           "test",
		Authentication: Authentication{},
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
		},
	},
}

var _ = Suite(dynamoTable_suite)

func (s *DynamoTableSuite) MapTest(c *C) {
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

	obtained := table.Map(tables["users_test"])

	c.Check(obtained, DeepEquals, expected)
}

func (s *DynamoTableSuite) TestGetByName(c *C) {
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
