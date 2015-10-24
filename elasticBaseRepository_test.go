package main

import (
	. "gopkg.in/check.v1"
)

type ElasticBaseRepositorySuite struct {
}

func (s *ElasticBaseRepositorySuite) SetUpTest(c *C) {
}

func (s *ElasticBaseRepositorySuite) TearDownTest(c *C) {
	var index ElasticIndex
	index.Delete("users_test")
}

var item = []Attribute{
	Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Monika"},
	Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Test"},
	Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "monika@gmail.com"},
	Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Poland"},
}

var elasticBaseRepository_suite = &ElasticBaseRepositorySuite{}

var _ = Suite(elasticBaseRepository_suite)

func (s *ElasticBaseRepositorySuite) TestAddItemByHash(c *C) {
	// Arrange
	var repo ElasticBaseRepository

	// Act
	obtained, _ := repo.Add("users_test", "id", "", item)

	// Assert
	c.Check(obtained, Equals, true)
}

func (s *ElasticBaseRepositorySuite) TestAddItemByHashRange(c *C) {
	// Arrange
	var repo ElasticBaseRepository

	// Act
	obtained, _ := repo.Add("users_test", "1", "Test", item)

	// Assert
	c.Check(obtained, Equals, true)
}

func (s *ElasticBaseRepositorySuite) TestDeleteItemByHash(c *C) {
	// Arrange
	var repo ElasticBaseRepository
	repo.Add("users_test", "1", "", item)

	// Act
	obtained, err := repo.DeleteByHash("users_test", "1")

	// Assert
	c.Check(obtained, Equals, true)
	c.Assert(err, IsNil)
}

func (s *ElasticBaseRepositorySuite) TestDeleteItemByHashRange(c *C) {
	// Arrange
	var repo ElasticBaseRepository
	repo.Add("users_test", "1", "test_range", item)

	// Act
	obtained, err := repo.DeleteByHashRange("users_test", "1", "test_range")

	// Assert
	c.Check(obtained, Equals, true)
	c.Assert(err, IsNil)
}
