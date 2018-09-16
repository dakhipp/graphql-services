//go:generate gorunpkg github.com/99designs/gqlgen --schema ../schema.graphql
package graph

import (
	"fmt"
	"log"
	"time"

	"github.com/dakhipp/graphql-services/auth"
	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

// global constants, might move to its own file at some time
// key type is needed to avoid warning: "should not use basic type string as key in context.withValue"
type key string

const CONTEXT_SESSION_KEY key = "user"
const CONTEXT_WRITER_KEY key = "writer"
const SESSION_COOKIE_NAME = "SESSION_ID"

// configuration struct created from environment variables
type envConfig struct {
	AuthURL    string `envconfig:"AUTH_SERVICE_URL"`
	RedisURL   string `envconfig:"REDIS_URL"`
	Port       string `envconfig:"PORT"`
	Playground bool   `envconfig:"PLAYGROUND"`
	Domain     string `envconfig:"DOMAIN"`
}

// GraphQLServer is a struct containing all server dependencies
type GraphQLServer struct {
	cfg             envConfig
	mux             *chi.Mux
	authClient      *auth.Client
	redisRepository Redis
}

// New registers all server dependencies needed for a GraphQLServer struct, registers routes and handles, and return a http ListenAndServe handler
func New() *GraphQLServer {
	// attempt to cast env variables into EnvConfig struct
	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// create new auth gRPC client
	authClient, err := auth.NewClient(cfg.AuthURL)
	if err != nil {
		fmt.Println("fail")
	}

	// create redis repository
	var redisRepository Redis
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		redisRepository, err = NewRedisRepository(cfg.RedisURL)
		if err != nil {
			log.Println(err)
		}
		return
	})

	// create GraphQLServer
	s := &GraphQLServer{
		cfg,
		chi.NewRouter(),
		authClient,
		redisRepository,
	}

	// register routes with created GraphQLServer
	RegisterRoutes(s)

	// return GraphQLServer to be used as a http ListenAndServe handler
	return s
}
