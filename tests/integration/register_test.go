package register

import (
	"log"
	"testing"

	. "github.com/franela/goblin"
	"github.com/kelseyhightower/envconfig"

	"github.com/vektah/gqlgen/client"
)

// Config : Configuration values created from environment variables
type Config struct {
	GQLEndpoint string `envconfig:"GQL_ENDPOINT"`
}

func TestRegister(t *testing.T) {
	g := Goblin(t)
	// Declare and attempt to cast config struct
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	c := client.New(cfg.GQLEndpoint)

	g.Describe("Register", func() {
		g.It("Should create a new user ", func() {
			var resp struct {
				Register struct {
					ID        string
					FirstName string
					LastName  string
				}
			}
			c.MustPost(`
				mutation {
					register(user: {
						firstName:"First",
						lastName:"Last"
					}) {
						id,
						firstName,
						lastName
					}
				}`, &resp)

			g.Assert(resp.Register.FirstName).Equal("First")
			g.Assert(resp.Register.LastName).Equal("Last")
		})
	})
}
