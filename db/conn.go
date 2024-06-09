package connections

import (
	configs "beli_mang/cfg"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgConn(config configs.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	// Construct the connection string using the provided configuration
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.DbUsername,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName,
	)

	// Parse the connection string to create a pgxpool configuration
	dbconfig, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pool config: %v", err)
	}

	// Set the connection pool configuration
	dbconfig.MaxConnLifetime = 1 * time.Hour
	dbconfig.MaxConnIdleTime = 30 * time.Second
	dbconfig.HealthCheckPeriod = 5 * time.Second
	dbconfig.MaxConns = 95 // Adjust based on your database and workload
	dbconfig.MinConns = 5  // Adjust based on your database and workload

	// Create a new connection pool with the configuration
	pool, err := pgxpool.NewWithConfig(ctx, dbconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create pool: %v", err)
	}

	// Test the connection before returning
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return pool, nil
}
