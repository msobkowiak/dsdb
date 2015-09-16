package main

import (
	"log"

	"github.com/olivere/elastic"
)

type ElasticIndex struct {
}

func (i ElasticIndex) Create(indexName string) error {
	client, err := elastic.NewClient()
	if err != nil {
		return err
	}

	// Check if index exists
	exists, err := client.IndexExists(indexName).Do()
	if err != nil {
		return err
	}

	if !exists {
		createIndex, err := client.CreateIndex(indexName).Do()
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			log.Println("Error on creating index")
		}

		table, err := GetTableDescription(indexName, schema.Tables)
		if err != nil {
			return err
		}
		if table.HasGeoPoint() {
			field, err := table.GetGeoPointName()
			if err != nil {
				return err
			} else {
				var mapper GeoPointMapper
				mapper.MapIndex(indexName, field, client)
			}
		}
	}

	return nil
}

func (i ElasticIndex) Delete(indexName string) error {
	client, err := elastic.NewClient()
	if err != nil {
		return err
	}

	_, err = client.DeleteIndex(indexName).Do()
	if err != nil {
		return err
	}

	return nil
}

func (i ElasticIndex) DeleteAll(schema DbDescription) {

	for tableName, _ := range schema.Tables {
		i.Delete(tableName)
	}
}
