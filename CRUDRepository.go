package main

type CRUDRepository interface {
	GetAll(tableName string) ([]map[string]string, error)
	GetByHash(tableName, hash string) (map[string]string, error)
	GetByHashRange(tableName, hashKey, rangeKey string) (map[string]string, error)
	DeleteByHash(tableName, hash string) (bool, error)
	DeleteByHashRange(tableName, hashKey, rangeKey string) (bool, error)
	Add(tableName, hashKey, rangeKey string, item []Attribute) (bool, error)
}
