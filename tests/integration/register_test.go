package register

import (
	"testing"

	. "github.com/franela/goblin"

	"github.com/vektah/gqlgen/client"
)

func TestRegister(t *testing.T) {
	g := Goblin(t)

	c := client.New("http://0.0.0.0:8000/graphql")

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
