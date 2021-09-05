package main

import (
	"fmt"
	"log"
	"os"
)

func ethEndpoint(network string) string {
	return fmt.Sprintf("wss://%s.infura.io/ws/v3/%s", network, os.Getenv("INFURA_PROJECT_ID"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func pgFromEnv() string {
	user := getEnv("PGUSER", "localhost")
	password := getEnv("PGPASSWORD", "root")
	host := getEnv("PGHOST", "localhost")
	port := getEnv("PGPORT", "5432")
	dbname := os.Getenv("PGDATABASE")

	// validation
	if dbname == "" {
		log.Fatalln("no database name exported: export PGDATABASE=xxxxx")
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	// const pguri = "postgresql://postgres:RDLPWbx5hM3ra@localhost:5432/crypto"
}
