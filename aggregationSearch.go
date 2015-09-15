package main

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/olivere/elastic"
)

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
			return false, errors.New("Not possible to calculate metric " + metric + " on non numeric filed")
		}
	}

	return true, nil
}