package datastore

import (
	"context"
	"fmt"
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDatabase(pool *dockertest.Pool) (*gorm.DB, *dockertest.Resource) {
	// Start PostgreSQL container
	postgresResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=admissions_test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start PostgreSQL container: %s", err)
		os.Exit(1)
	}

	// Exponential backoff-retry to connect to PostgreSQL
	var db *gorm.DB
	err = pool.Retry(func() error {
		dsn := fmt.Sprintf("host=localhost port=%s user=postgres password=postgres dbname=admissions_test sslmode=disable",
			postgresResource.GetPort("5432/tcp"))
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		return db.Exec("SELECT 1").Error
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to PostgreSQL: %s", err)
		os.Exit(1)
	}

	return db, postgresResource
}

func setupCache(pool *dockertest.Pool) (*redis.Client, *dockertest.Resource) {
	// Start Redis container
	redisResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "latest",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start Redis container: %s", err)
		os.Exit(1)
	}

	// Exponential backoff-retry to connect to Redis
	var cache *redis.Client
	err = pool.Retry(func() error {
		cache = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp")),
		})
		_, err := cache.Ping(context.Background()).Result()
		return err
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to Redis: %s", err)
		os.Exit(1)
	}

	return cache, redisResource
}

func SetupMockStorage() (storage Storage, cleanup func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to docker: %s", err)
		os.Exit(1)
	}

	db, postgresResource := setupDatabase(pool)
	cache, redisResource := setupCache(pool)

	storage = NewStorage(db, cache)

	cleanup = func() {
		if err := pool.Purge(postgresResource); err != nil {
			fmt.Fprintf(os.Stderr, "Could not purge PostgreSQL container: %s", err)
		}
		if err := pool.Purge(redisResource); err != nil {
			fmt.Fprintf(os.Stderr, "Could not purge Redis container: %s", err)
		}
	}

	return
}
