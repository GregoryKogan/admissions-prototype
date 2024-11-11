package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPostgres(t *testing.T) *gorm.DB {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	// Start a PostgreSQL container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=admissions_test",
		},
	}, func(config *docker.HostConfig) {
		// Expose the container to the host
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(t, err)

	// Clean up the container after tests
	t.Cleanup(func() {
		err := pool.Purge(resource)
		require.NoError(t, err)
	})

	// Exponential backoff-retry to connect to the database
	var db *gorm.DB
	err = pool.Retry(func() error {
		dsn := fmt.Sprintf("host=localhost port=%s user=postgres password=postgres dbname=admissions_test sslmode=disable",
			resource.GetPort("5432/tcp"))
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		return db.
			Exec("SELECT 1").
			Error
	})
	require.NoError(t, err)

	require.NotNil(t, db)

	return db
}

func setupRedis(t *testing.T) *redis.Client {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	// Start a Redis container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "latest",
	}, func(config *docker.HostConfig) {
		// Expose the container to the host
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(t, err)

	// Clean up the container after tests
	t.Cleanup(func() {
		err := pool.Purge(resource)
		require.NoError(t, err)
	})

	// Exponential backoff-retry to connect to the cache
	var client *redis.Client
	err = pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})
		_, err := client.Ping(context.Background()).Result()
		return err
	})
	require.NoError(t, err)

	require.NotNil(t, client)

	return client
}

func SetupMockStorage(t *testing.T) Storage {
	db := setupPostgres(t)
	cache := setupRedis(t)
	return NewStorage(db, cache)
}
