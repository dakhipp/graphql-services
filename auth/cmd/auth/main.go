package main

import (
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	Port        string `envconfig:"PORT"`
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
		repository, err = auth.NewPostgresRepository(cfg.DatabaseURL)
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
