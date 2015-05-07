package main

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/dynamodb"
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
