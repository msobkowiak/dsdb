package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func LoadSchema(filePath string) DbDescription {
	body, err := ioutil.ReadFile(filePath)
	check(err)

	var configData struct {
		Database_name     string
		Dynamo_region     string
		Dynamo_access_key string
		Dynamo_secret_key string
		Tables            []struct {
			Name              string
			Attributes        []AttributeDefinition
			Primary_key       KeyDefinition
			Secondary_indexes []SecondaryIndexDefinition
			Full_text_search  []TextSearch
		}
	}

	err = yaml.Unmarshal([]byte(body), &configData)
	if err != nil {
		log.Println(err)
	}

	t := make(map[string]TableDescription, len(configData.Tables))
	for _, table := range configData.Tables {
		t[table.Name] = TableDescription{
			Name:             table.Name,
			Attributes:       table.Attributes,
			PrimaryKey:       table.Primary_key,
			SecondaryIndexes: table.Secondary_indexes,
		}
	}

	schema := DbDescription{
		Name: configData.Database_name,
		Authentication: Authentication{
			Dynamo: DynamoAuth{
				Region:    configData.Dynamo_region,
				AccessKey: configData.Dynamo_access_key,
				SecretKey: configData.Dynamo_secret_key,
			},
		},
		Tables: t,
	}

	return schema
}
