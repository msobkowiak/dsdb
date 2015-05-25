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

	var data = make(map[string]string)
	for i := range item {
		data[item[i].Description.Name] = item[i].Value
	}
	if rangeValue != "" {
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
}