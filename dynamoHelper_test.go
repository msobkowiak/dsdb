package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
	"reflect"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func SetUpSuite(c *C) {
}

func (s *MySuite) TestAuth(c *C) {
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

func (s *MySuite) TestConvertToDynamo(c *C) {
	var table = TableDescription{
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
	}

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
				IndexName: "counrty",
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

	obtained := ConvertToDynamo(table)

	reflect.DeepEqual(expected, obtained)
}
