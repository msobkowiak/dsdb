package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/olivere/elastic"
)

func GeoSearch(tableName, field, distance string, lat, lon float64) ([]map[string]string, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	allQuery := elastic.NewMatchAllQuery()
	builder := client.Search().
		Index(tableName).
		Query(&allQuery)

	f := elastic.NewGeoDistanceFilter(field)
	f = f.Lat(lat)
	f = f.Lon(lon)
	f = f.Distance(distance)
	f = f.DistanceType("plane")
	f = f.OptimizeBbox("memory")

	builder.PostFilter(f)

	searchResult, err := builder.Do()
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
