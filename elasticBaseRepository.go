package main

import (
	"encoding/json"
	"log"

	"github.com/olivere/elastic"
)

type ElasticBaseRepository struct {
}

func (r ElasticBaseRepository) Add(tableName, hashKey, rangeKey string, item []Attribute) (bool, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return false, err
	}

	var index ElasticIndex
	index.Create(tableName)

	doc, data := r.mapToDocument(tableName, hashKey, rangeKey, item)

	indexBody, _ := json.Marshal(data)

	_, err = client.Index().
		Index(doc["indexName"]).
		Id(doc["id"]).
		Type(doc["indexType"]).
		BodyJson(string(indexBody)).
		Do()
	if err != nil {
		return false, err
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index(tableName).Do()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r ElasticBaseRepository) DeleteByHash(tableName, hash string) (bool, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return false, err
	}

	_, err = client.Delete().
		Index(tableName).
		Type(tableName).
		Id(hash).
		Do()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r ElasticBaseRepository) DeleteByHashRange(tableName, hashKey, rangeKey string) (bool, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return false, err
	}

	_, err = client.Delete().
		Index(tableName).
		Type(tableName).
		Id(hashKey + "_" + rangeKey).
		Do()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r ElasticBaseRepository) mapToDocument(tableName, hashKey, rangeKey string, item []Attribute) (map[string]string, map[string]interface{}) {
	data := map[string]interface{}{}

	for i := range item {
		if item[i].Description.Type == "G" {
			var mapper GeoPointMapper
			geo, _ := mapper.MapStringToGeoPoint(item[i].Value.(string))
			data[item[i].Description.Name] = geo
		} else {
			data[item[i].Description.Name] = item[i].Value
		}
	}

	if rangeKey != "" {
		hashName, err := schema.GetHashName(tableName)
		if err != nil {
			log.Println(err)
		}

		rangeName, err := schema.GetRangeName(tableName)
		if err != nil {
			log.Println(err)
		}

		data[hashName] = hashKey
		data[rangeName] = rangeKey
		hashKey = hashKey + "_" + rangeKey
	}

	indexDoc := map[string]string{
		"indexName": tableName,
		"indexType": tableName,
		"id":        hashKey,
	}

	return indexDoc, data
}
