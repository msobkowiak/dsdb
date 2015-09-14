package main

import (
	"encoding/json"
	"log"

	"github.com/olivere/elastic"
)

func AddToElasticSearch(indexName, indexType, idValue, rangeValue string, item []Attribute) {
	client, err := elastic.NewClient()
	if err != nil {
		log.Println(err)
	}

	createIndex(indexName, client)

	data := map[string]interface{}{}
	for i := range item {
		if item[i].Description.Name == "location" {
			geo, _ := GeoPointFromString(item[i].Value.(string))
			data[item[i].Description.Name] = geo
		} else {
			data[item[i].Description.Name] = item[i].Value
		}
	}

	if rangeValue != "" {
		hashName, err := GetHashName(indexType, schema)
		if err != nil {
			log.Println(err)
		}

		rangeName, err := GetRangeName(indexType, schema)
		if err != nil {
			log.Println(err)
		}

		data[hashName] = idValue
		data[rangeName] = rangeValue
		idValue = idValue + "_" + rangeValue
	}

	indexBody, _ := json.Marshal(data)
	addIndexValue(indexName, indexType, idValue, indexBody, client)
}

func createIndex(indexName string, client *elastic.Client) {
	// Check if index exists
	exists, err := client.IndexExists(indexName).Do()
	if err != nil {
		log.Println(err)
	}

	if !exists {
		createIndex, err := client.CreateIndex(indexName).Do()
		if err != nil {
			log.Println(err)
		}
		if !createIndex.Acknowledged {
			log.Println("Error on creating index")
		}

		table, err := GetTableDescription(indexName, schema.Tables)
		if err != nil {
			log.Println(err)
		}
		if table.HasGeoPoint() {
			field, err := table.GetGeoPointName()
			if err != nil {
				log.Println(err)
			} else {
				MappGeoPoint(indexName, field, client)
			}
		}
	}
}

func addIndexValue(indexName, indexType, id string, indexBody []byte, client *elastic.Client) {

	_, err := client.Index().
		Index(indexName).
		Id(id).
		Type(indexType).
		BodyJson(string(indexBody)).
		Do()
	if err != nil {
		log.Println(err)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index(indexName).Do()
	if err != nil {
		panic(err)
	}
}
