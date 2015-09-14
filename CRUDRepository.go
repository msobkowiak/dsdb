package main

type CRUDRepository interface {
	GetAll(tableName string) ([]map[string]string, error)
	GetByHash(tableName, hash string) (map[string]string, error)
	GetHashRange(tableName, hashKey, rangeKey string) (map[string]string, error)
	//GetItemsByHash(tableName, hashValue string) ([]map[string]string, error)
	//GetByIndexHash(tableName, indexName, hashValue string) ([]map[string]string, error)
	//GetItemsByRangeOp(tableName, hashValue, operator string, rangeValue []string) ([]map[string]string, error)
	//GetByIndexRangeOp(tableName, indexName, hashValue, operator string, rangeValue []string) ([]map[string]string, error)

	DeleteByHash(tableName, hash string) (bool, error)
	DeleteByHashRange(tableName, hashKey, rangeKey string) (bool, error)

	Add(tableName, hashKey, rangeKey string, item []Attribute) (bool, error)
}
