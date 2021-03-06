package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/olivere/elastic"
)

type GeoPointMapper struct {
}

type geoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (m *GeoPointMapper) MapToIndex(indexName, field string, client *elastic.Client) error {
	mapping := `{
		"` + indexName + `":{
			"properties":{
				"` + field + `":{
					"type":"geo_point"
				}
			}
		}
	}`

	_, err := client.PutMapping().Index(indexName).Type(indexName).BodyString(mapping).Do()
	if err != nil {
		return err
	}

	return nil
}

func (m *GeoPointMapper) MapStringToGeoPoint(latLon string) (geoPoint, error) {
	latlon := strings.SplitN(latLon, ",", 2)
	if len(latlon) != 2 {
		return geoPoint{}, errors.New(latLon + " is not a valid geo point string")
	}
	latValue, err := strconv.ParseFloat(latlon[0], 64)
	if err != nil {
		return geoPoint{}, err
	}
	lonValue, err := strconv.ParseFloat(latlon[1], 64)
	if err != nil {
		return geoPoint{}, err
	}
	return geoPoint{Lat: latValue, Lon: lonValue}, nil
}
