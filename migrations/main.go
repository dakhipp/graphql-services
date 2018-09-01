package main

import (
	"log"
	"os"
	"time"

	"github.com/go-pg/pg"
	"github.com/kelseyhightower/envconfig"
	migrations "github.com/robinjoseph08/go-pg-migrations"
	"github.com/tinrab/retry"
)

// Config : Configuration values created from environment variables
type Config struct {
	PSQLUser string `envconfig:"PSQL_USER"`
	PSQLPass string `envconfig:"PSQL_PASS"`
	PSQLAddr string `envconfig:"PSQL_ADDR"`
	PSQLDB   string `envconfig:"PSQL_DB"`
}

const directory = "."

func main() {
	// Declare and attempt to cast config struct
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		db := pg.Connect(&pg.Options{
			User:     cfg.PSQLUser,
			Password: cfg.PSQLPass,
			Addr:     cfg.PSQLAddr,
			Database: cfg.PSQLDB,
		})

		err = migrations.Run(db, directory, os.Args)
		if err != nil {
			log.Println(err)
		}
		return
	})
}
