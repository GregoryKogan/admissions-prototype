package datastore

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Storage interface {
	DB() *gorm.DB
	Cache() *redis.Client
}

type StorageImpl struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewStorage(db *gorm.DB, cache *redis.Client) Storage {
	return StorageImpl{
		db:    db,
		cache: cache,
	}
}

func InitStorage() Storage {
	db := InitDBConnection()
	cache := InitCacheConnection()
	return NewStorage(db, cache)
}

func (s StorageImpl) DB() *gorm.DB {
	return s.db
}

func (s StorageImpl) Cache() *redis.Client {
	return s.cache
}
