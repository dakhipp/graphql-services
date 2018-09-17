package main

import (
	"fmt"
	"log"

	"github.com/dakhipp/graphql-services/auth"
	"github.com/kelseyhightower/envconfig"
)

// configuration struct created from environment variables
type envConfig struct {
	Port string `envconfig:"PORT"`
}

func main() {
	// attempt to cast env variables into envConfig struct
	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// log when the server starts
	log.Println(fmt.Sprintf("Listening on port %s...", cfg.Port))

	// start GRPC server or throw fatal error
	log.Fatal(auth.ListenGRPC(auth.New(), fmt.Sprintf(":%s", cfg.Port)))
}
