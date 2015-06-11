package main

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/olivere/elastic"

	"math"
)

func AddToElasticSearch(indexName, indexType, idValue, rangeValue string, item []Attribute) {
	client, err := elastic.NewClient()
	if err != nil {
		log.Println(err)
	}

	createIndex(indexName, client)

	data := map[string]interface{}{}
	for i := range item {
		if item[i].Description.Type == "S" {
			data[item[i].Description.Name] = item[i].Value
		} else {
			value, _ := strconv.Atoi(item[i].Value)
			data[item[i].Description.Name] = value
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

	var res struct {
		Value float64
	}
	err = json.Unmarshal(*aggResult, &res)
	if err != nil {
		return nil, err
	}

	var result = make(map[string]float64, 1)
	result[metric] = round(res.Value, 2)

	return result, nil
}

func AggregationStatsSearch(index, field, metric string) (map[string]float64, error) {
	log.Println(metric)
	client, err := elastic.NewClient()
	if err != nil {
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

	var stats = make(map[string]float64, 8)
	json.Unmarshal(*aggResult, &stats)

	for res := range stats {
		stats[res] = round(stats[res], 2)
	}
	return stats, nil
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
