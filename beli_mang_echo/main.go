package main

import (
	"log"

	configs "beli_mang/cfg"
	connections "beli_mang/db"
	"beli_mang/server"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Set up database connection pool
	dbPool, err := connections.NewPgConn(config)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbPool.Close()

	s := server.NewServer(dbPool)

	s.RegisterRoute(config)

	s.Run(config)
}
