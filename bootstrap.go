package main

import (
	"strconv"
)

var data = map[string][][]Attribute{
	"users": [][]Attribute{
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Monika"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Sobkowiak"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "monika@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Poland"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Ana"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Dias"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "ana@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Nuno"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Correia"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "nuno@example.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Nuno"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Correia"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "nuno@example.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Isabel"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Fernendes"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "isabel@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Spain"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Miguel"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Oliveira"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "miguel@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Mikolaj"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Nowak"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "mikolaj@example.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Poland"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "first_name", Type: "S"}, Value: "Joao"},
			Attribute{Description: AttributeDefinition{Name: "last_name", Type: "S"}, Value: "Silva"},
			Attribute{Description: AttributeDefinition{Name: "email", Type: "S"}, Value: "joao@gmail.com"},
			Attribute{Description: AttributeDefinition{Name: "country", Type: "S"}, Value: "Portugal"},
		},
	},
	"game_scores": [][]Attribute{
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: 123},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: 20},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: 0},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: 333},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: 90},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: 21},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: 444},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: 99},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: 59},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: 555},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: 12},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: 9},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: 666},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: 7},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: 20},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: 122},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: 20},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: 7},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: 777},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: 190},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: 87},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "top_score", Type: "N"}, Value: 888},
			Attribute{Description: AttributeDefinition{Name: "wins", Type: "N"}, Value: 438},
			Attribute{Description: AttributeDefinition{Name: "losts", Type: "N"}, Value: 164},
		},
	},
	"restaurants": [][]Attribute{
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "name", Type: "S"}, Value: "Portuguese Food"},
			Attribute{Description: AttributeDefinition{Name: "location", Type: "G"}, Value: "41.158915,-8.6191053"},
			Attribute{Description: AttributeDefinition{Name: "descripion", Type: "S"}, Value: "The best tradicional portugues restaurant in town. Our speciality is cotfish!"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "name", Type: "S"}, Value: "Italian Food"},
			Attribute{Description: AttributeDefinition{Name: "location", Type: "G"}, Value: "41.158915,-8.631091"},
			Attribute{Description: AttributeDefinition{Name: "descripion", Type: "S"}, Value: "The best tradicional italian passta restaurant in town."},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "name", Type: "S"}, Value: "Pasta Cafe"},
			Attribute{Description: AttributeDefinition{Name: "location", Type: "G"}, Value: "41.159238,-8.621091"},
			Attribute{Description: AttributeDefinition{Name: "descripion", Type: "S"}, Value: "The best tradicional portugues restaurant in town. Our speciality: pasta, pasta, pasta!"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "name", Type: "S"}, Value: "Pasta Cafe1"},
			Attribute{Description: AttributeDefinition{Name: "location", Type: "G"}, Value: "41.154941,-8.638"},
			Attribute{Description: AttributeDefinition{Name: "descripion", Type: "S"}, Value: "The best tradicional portugues restaurant in town. Our speciality: pasta, pasta, pasta!"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "name", Type: "S"}, Value: "Pasta Cafe2"},
			Attribute{Description: AttributeDefinition{Name: "location", Type: "G"}, Value: "41.15002,-8.635297"},
			Attribute{Description: AttributeDefinition{Name: "descripion", Type: "S"}, Value: "The best tradicional portugues restaurant in town. Our speciality: pasta, pasta, pasta!"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "name", Type: "S"}, Value: "Pasta Cafe3"},
			Attribute{Description: AttributeDefinition{Name: "location", Type: "G"}, Value: "41.147832,-8.611693"},
			Attribute{Description: AttributeDefinition{Name: "descripion", Type: "S"}, Value: "The best tradicional portugues restaurant in town. Our speciality: pasta, pasta, pasta!"},
		},
		[]Attribute{
			Attribute{Description: AttributeDefinition{Name: "name", Type: "S"}, Value: "Pasta Cafe4"},
			Attribute{Description: AttributeDefinition{Name: "location", Type: "G"}, Value: "41.147379,-8.605642"},
			Attribute{Description: AttributeDefinition{Name: "descripion", Type: "S"}, Value: "The best tradicional portugues restaurant in town. Our speciality: pasta, pasta, pasta!"},
		},
	},
}

var hashKeys = map[string][]string{
	"game_scores": []string{
		"Game Brown Fox",
		"Game Green Fox",
		"Game X",
		"Game Brown",
		"Game Green",
		"Game Fox",
		"Game X",
		"Game Y",
	},
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Bootstrap(dbDescription DbDescription) {

	var (
		dynamoRepo   DynamoBaseRepository
		dynamoClient DynamoClient
		dynamoTable  DynamoTable
		elasticRepo  ElasticBaseepository
		elasticIndex ElasticIndex
	)

	auth := dbDescription.Authentication.Dynamo
	db := dynamoClient.Auth(auth.Region, auth.AccessKey, auth.SecretKey)

	// cleanup the databases
	dynamoClient.DeleteAll(db)
	elasticIndex.DeleteAll(dbDescription)

	//create tables with example data
	tableName := "users"
	dynamoTable.Create(dbDescription.Tables[tableName])
	for i := range data[tableName] {
		hash := strconv.FormatInt(int64(i+1), 10)
		dynamoRepo.Add(tableName, hash, "", data[tableName][i])
		elasticRepo.Add(tableName, hash, "", data[tableName][i])
	}

	tableName = "game_scores"
	dynamoTable.Create(dbDescription.Tables[tableName])
	for i := range data[tableName] {
		rangeValue := strconv.FormatInt(int64(i+1), 10)
		dynamoRepo.Add(tableName, hashKeys[tableName][i], rangeValue, data[tableName][i])
		elasticRepo.Add(tableName, hashKeys[tableName][i], rangeValue, data[tableName][i])
	}

	tableName = "restaurants"
	dynamoTable.Create(dbDescription.Tables[tableName])
	for i := range data[tableName] {
		hash := strconv.FormatInt(int64(i+1), 10)
		dynamoRepo.Add(tableName, hash, "", data[tableName][i])
		elasticRepo.Add(tableName, hash, "", data[tableName][i])
	}
}
