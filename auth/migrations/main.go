package main

import (
	"log"
	"os"

	migrations "github.com/dakhipp/go-pg-migrations"
	"github.com/go-pg/pg"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PSQLUser string `envconfig:"PSQL_USER"`
	PSQLPass string `envconfig:"PSQL_PASS"`
	PSQLAddr string `envconfig:"PSQL_ADDR"`
	PSQLDB   string `envconfig:"PSQL_DB"`
}

const directory = "auth/migrations"

func main() {
	// Declare and attempt to cast config struct
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	db := pg.Connect(&pg.Options{
		User:     cfg.PSQLUser,
		Password: cfg.PSQLPass,
		Addr:     cfg.PSQLAddr,
		Database: cfg.PSQLDB,
	})

	// db := pg.Connect(&pg.Options{
	// 	User:     "root",
	// 	Password: "toor",
	// 	Addr:     "localhost:5432",
	// 	Database: "psql",
	// })

	err = migrations.Run(db, directory, os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
