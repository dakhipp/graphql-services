package main

import (
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	Port     string `envconfig:"PORT"`
	PSQLUser string `envconfig:"PSQL_USER"`
	PSQLPass string `envconfig:"PSQL_PASS"`
	PSQLAddr string `envconfig:"PSQL_ADDR"`
	PSQLDB   string `envconfig:"PSQL_DB"`
	PSQLSSL  string `envconfig:"PSQL_SSL"`
}

func main() {
	// Declare and attempt to cast config struct
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Attempt to create auth repository
	var repository auth.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		databaseURL := "postgres://" + cfg.PSQLUser + ":" + cfg.PSQLPass + "@" + cfg.PSQLAddr + "/" + cfg.PSQLDB + "?sslmode=" + cfg.PSQLSSL
		repository, err = auth.NewPostgresRepository(databaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()

	// Log when the server starts
	log.Println("Listening on port " + cfg.Port + "...")

	// Create service from repository
	service := auth.NewService(repository)

	// Start or throw fatal error
	log.Fatal(auth.ListenGRPC(service, ":"+cfg.Port))
}
