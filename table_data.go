package main

import (
	"github.com/goamz/goamz/dynamodb"
)

// structure of users table
var users = dynamodb.TableDescriptionT{
	TableName: "Users",
	AttributeDefinitions: []dynamodb.AttributeDefinitionT{
		dynamodb.AttributeDefinitionT{"id", "N"},
	},
	KeySchema: []dynamodb.KeySchemaT{
		dynamodb.KeySchemaT{"id", "HASH"},
	},
	ProvisionedThroughput: dynamodb.ProvisionedThroughputT{
		ReadCapacityUnits:  40,
		WriteCapacityUnits: 20,
	},
}

func LoadUsersData() [][]dynamodb.Attribute {
	var users_data = make([][]dynamodb.Attribute, 8)
	users_data[0] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Monika"),
		*dynamodb.NewStringAttribute("last_name", "Sobkowiak"),
		*dynamodb.NewStringAttribute("email", "monika@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Poland"),
	}
	users_data[1] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Ana"),
		*dynamodb.NewStringAttribute("last_name", "Dias"),
		*dynamodb.NewStringAttribute("email", "ana@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Portugal"),
	}
	users_data[2] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Nuno"),
		*dynamodb.NewStringAttribute("last_name", "Correia"),
		*dynamodb.NewStringAttribute("email", "nuno@exemple.com"),
		*dynamodb.NewStringAttribute("counrty", "Portugal"),
	}
	users_data[3] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Isabel"),
		*dynamodb.NewStringAttribute("last_name", "Frenandes"),
		*dynamodb.NewStringAttribute("email", "isabel@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Spain"),
	}
	users_data[4] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Miguel"),
		*dynamodb.NewStringAttribute("last_name", "Oliveira"),
		*dynamodb.NewStringAttribute("email", "miguel@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Portugal"),
	}
	users_data[5] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Mikolaj"),
		*dynamodb.NewStringAttribute("last_name", "Nowak"),
		*dynamodb.NewStringAttribute("email", "mikolaj@exemple.com"),
		*dynamodb.NewStringAttribute("counrty", "Poland"),
	}
	users_data[6] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Joao"),
		*dynamodb.NewStringAttribute("last_name", "Silva"),
		*dynamodb.NewStringAttribute("email", "joao@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "Portugal"),
	}
	users_data[7] = []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("first_name", "Mat"),
		*dynamodb.NewStringAttribute("last_name", "Deamon"),
		*dynamodb.NewStringAttribute("email", "mat@gmail.com"),
		*dynamodb.NewStringAttribute("counrty", "USA"),
	}

	return users_data
}
