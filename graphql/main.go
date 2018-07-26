package main

import (
	"log"
	"net/http"

	"github.com/dakhipp/graphql-services/graphql/graph"
	"github.com/kelseyhightower/envconfig"
	"github.com/vektah/gqlgen/handler"
)

type Config struct {
	AuthUrl string `envconfig:"AUTH_SERVICE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := graph.NewGraphQLServer(cfg.AuthUrl)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/graphql", handler.GraphQL(graph.MakeExecutableSchema(s)))
	http.Handle("/playground", handler.Playground("Auth", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
