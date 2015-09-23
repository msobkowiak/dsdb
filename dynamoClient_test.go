package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"

	. "gopkg.in/check.v1"
)

func (s *TableSuite) TestAuth(c *C) {
	var client DynamoClient
	region := "http://127.0.0.1:4567"
	accessKey := "key"
	secretKey := "secret"
	obtained := client.Auth(region, accessKey, secretKey)

	expected := dynamodb.Server{
		Auth:   aws.Auth{AccessKey: "key", SecretKey: "secret"},
		Region: aws.Region{DynamoDBEndpoint: "http://127.0.0.1:4567"},
	}

	c.Check(obtained, Equals, expected)
}
