package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/olivere/elastic"

	"math"
)

type geoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

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
				mappGeoPoint(indexName, field, client)
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

func mappGeoPoint(indexName, field string, client *elastic.Client) {
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

func AggregationSearch(index, field, metric string) (map[string]float64, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	ok, err := isCalculable(index, metric, field)
	if !ok {
		return nil, err
	}

	allQuery := elastic.NewMatchAllQuery()
	builder := client.Search().
		Index(index).
		Query(&allQuery)

	builder, err = decorate(metric, field, builder)
	if err != nil {
		return nil, err
	}

	searchResult, _ := builder.Do()
	aggResult, _ := searchResult.Aggregations[metric]
	if err != nil {
		return nil, err
	}

	var result = map[string]float64{}
	json.Unmarshal(*aggResult, &result)

	if metric == "stats" || metric == "extended_stats" {
		for i := range result {
			result[i] = round(result[i], 2)
		}
	} else {
		var res = map[string]float64{}
		res[metric] = round(result["value"], 2)
		result = res
	}

	return result, nil
}

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

func decorate(metric, field string, builder *elastic.SearchService) (*elastic.SearchService, error) {
	switch metric {
	case "max":
		return max(field, builder), nil
	case "min":
		return min(field, builder), nil
	case "sum":
		return sum(field, builder), nil
	case "avg":
		return avg(field, builder), nil
	case "count":
		return count(field, builder), nil
	case "stats":
		return stats(field, builder), nil
	case "extended_stats":
		return extendedStats(field, builder), nil
	}

	return nil, errors.New("Metric " + metric + " not known.")
}

func max(field string, builder *elastic.SearchService) *elastic.SearchService {
	agg := elastic.NewMaxAggregation().Field(field)
	return builder.Aggregation("max", agg.Field(field))
}

func min(field string, builder *elastic.SearchService) *elastic.SearchService {
	agg := elastic.NewMinAggregation().Field(field)
	return builder.Aggregation("min", agg.Field(field))
}

func sum(field string, builder *elastic.SearchService) *elastic.SearchService {
	agg := elastic.NewSumAggregation().Field(field)
	return builder.Aggregation("sum", agg.Field(field))
}

func avg(field string, builder *elastic.SearchService) *elastic.SearchService {
	agg := elastic.NewAvgAggregation().Field(field)
	return builder.Aggregation("avg", agg.Field(field))
}

func count(field string, builder *elastic.SearchService) *elastic.SearchService {
	agg := elastic.NewValueCountAggregation().Field(field)
	return builder.Aggregation("count", agg.Field(field))
}

func stats(field string, builder *elastic.SearchService) *elastic.SearchService {
	agg := elastic.NewStatsAggregation().Field(field)
	return builder.Aggregation("stats", agg.Field(field))
}

func extendedStats(field string, builder *elastic.SearchService) *elastic.SearchService {
	agg := elastic.NewExtendedStatsAggregation().Field(field)
	return builder.Aggregation("extended_stats", agg.Field(field))
}

func round(val float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= .5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	return round / pow
}

func isCalculable(index, metric, field string) (bool, error) {
	if metric != "count" {
		table, err := GetTableDescription(index, schema.Tables)
		if err != nil {
			return false, err
		}

		metricType := table.GetTypeOfAttribute(field)
		if metricType != "N" {
			return false, errors.New("Not possible to calculate metric " + metric + " non numeric filed")
		}
	}

	return true, nil
}

func GeoPointFromString(latLon string) (geoPoint, error) {
	latlon := strings.SplitN(latLon, ",", 2)
	if len(latlon) != 2 {
		return geoPoint{}, errors.New("elastic: " + latLon + " is not a valid geo point string")
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
