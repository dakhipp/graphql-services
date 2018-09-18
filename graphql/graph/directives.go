package graph

import (
	"context"
	"fmt"

	graphql "github.com/99designs/gqlgen/graphql"
)

func (s *GraphQLServer) hasRole(ctx context.Context, next graphql.Resolver, role Role) (interface{}, error) {
	// pull current session off of context
	a, ok := ctx.Value(CONTEXT_SESSION_KEY).(Session)

	if ok {
		for _, v := range a.Roles {
			if v == role {
				return next(ctx)
			}
		}
	}

	return nil, fmt.Errorf("Access denied")
}
