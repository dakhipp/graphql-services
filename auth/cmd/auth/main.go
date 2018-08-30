package main

import (
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

// Config : Configuration values created from environment variables
type Config struct {
	Port    string `envconfig:"PORT"`
	PSQLURL string `envconfig:"PSQL_URL"`
}

func main() {
	// Declare and attempt to cast config struct
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg.PSQLURL)

	// Attempt to create auth repository
	var repository auth.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = auth.NewPostgresRepository(cfg.PSQLURL)
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
