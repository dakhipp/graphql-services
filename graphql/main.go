package main

import (
	"log"
	"net/http"

	"github.com/dakhipp/graphql-services/graphql/graph"
	"github.com/kelseyhightower/envconfig"
	"github.com/vektah/gqlgen/handler"
)

// Config is an exported struct containing environment variables.
type Config struct {
	AuthURL    string `envconfig:"AUTH_SERVICE_URL"`
	Port       string `envconfig:"PORT"`
	Playground bool   `envconfig:"PLAYGROUND"`
}

func main() {
	// Declare and attempt to cast config struct
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Can take multiple comma separated config urls
	server, err := graph.NewGraphQLServer(cfg.AuthURL)
	if err != nil {
		log.Fatal(err)
	}

	// Register graphql route
	http.Handle("/graphql", handler.GraphQL(graph.NewExecutableSchema(graph.Config{Resolvers: server})))

	// Register playgorund route if environment variable is set to true
	if cfg.Playground == true {
		http.Handle("/playground", handler.Playground("Graphql Playground", "/graphql"))
	}

	// Log when the server starts
	log.Println("Listening on port " + cfg.Port + "...")

	// Start or throw fatal error
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
