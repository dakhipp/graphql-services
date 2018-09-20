package graph

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/vektah/gqlgen/handler"
)

// RegisterRoutes takes a GraphQLServer pointer and registers all needed middlewares, router handlers, and directives
func RegisterRoutes(s *GraphQLServer) {
	gqlConfig := Config{
		Resolvers: s,
	}

	// A good base middleware stack
	s.mux.Use(middleware.RequestID)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)

	// register global middleware to attach a http.ResponseWriter to the context of the request
	s.mux.Use(s.attachResponseWriterMiddleware())

	// register global middleware to attach a user based on the session cookie to the context of the request
	s.mux.Use(s.attachUserMiddleware())

	// register directives
	gqlConfig.Directives.HasRole = s.hasRole

	// register GraphQL route
	s.mux.Handle("/graphql", handler.GraphQL(NewExecutableSchema(gqlConfig)))

	// register healthcheck route
	s.mux.HandleFunc("/h", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "healthy") })

	// register playground route if environment variable is set to true
	if s.cfg.Playground {
		s.mux.Handle("/playground", handler.Playground("Graphql Playground", "/graphql"))
	}
}

// allows us to use GraphQLServer with a mux attached to it as a router
func (s *GraphQLServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// attachResponseWriter is middleware used to attach the http.ResponseWriter to context in order to send a cookie from inside resolver functions
func (s *GraphQLServer) attachResponseWriterMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// put ResponseWriter into context
			ctx := context.WithValue(r.Context(), contextWriterKey, w)

			// call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// attachUserMiddleware validates the session ID cookie and attaches a valid session from redis to the request context
func (s *GraphQLServer) attachUserMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get session cookie
			c, err := r.Cookie(sessionCookieName)

			// allow unauthenticated users in
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			// get session ID and attempt to look up the session in redis
			sID := c.Value
			ses, err := s.redisRepository.GetSession(sID)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(fmt.Sprintf("Currently logged in user: %v", ses))

			// attach session to context, it might be an empty session if there is an error
			ctx := context.WithValue(r.Context(), contextSessionKey, ses)

			// attach session ID to context, this is needed in order to delete the session
			ctx2 := context.WithValue(ctx, contextSessionIDKey, sID)

			// call the next with our new context
			r = r.WithContext(ctx)
			r = r.WithContext(ctx2)
			next.ServeHTTP(w, r)
		})
	}
}
