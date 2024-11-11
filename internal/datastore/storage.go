package datastore

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Storage struct {
	DB    *gorm.DB
	Cache *redis.Client
}

func NewStorage(db *gorm.DB, cache *redis.Client) Storage {
	return Storage{
		DB:    db,
		Cache: cache,
	}
}

func InitStorage() Storage {
	db := InitDBConnection()
	cache := InitCacheConnection()
	return NewStorage(db, cache)
}
