package main

import (
	"encoding/json"
	"errors"
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
		data["game_title"] = idValue
		data["user_id"] = rangeValue
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

func FullTextSearchQuery(index, field, query, operator, precision string) ([]map[string]string, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	matchQuery := elastic.NewMatchQuery(field, query).
		Operator(operator)

	if precision != "" {
		matchQuery.MinimumShouldMatch(precision + "%")
	}

	searchResult, err := client.Search().
		Index(index).
		Query(&matchQuery).
		Pretty(true).
		Do()
	if err != nil {
		return nil, err
	}

	if searchResult.Hits != nil {
		var result = make([]map[string]string, searchResult.Hits.TotalHits)
		for i, hit := range searchResult.Hits.Hits {
			var t map[string]string
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				log.Println(err)
			}
			result[i] = t
		}
		return result, nil
	}

	return nil, errors.New("No results found")
}
