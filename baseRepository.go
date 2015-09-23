package main

type BaseRepository interface {
	DeleteByHash(tableName, hash string) (bool, error)
	DeleteByHashRange(tableName, hashKey, rangeKey string) (bool, error)
	Add(tableName, hashKey, rangeKey string, item []Attribute) (bool, error)
}
