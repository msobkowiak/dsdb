package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/olivere/elastic"
)

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
