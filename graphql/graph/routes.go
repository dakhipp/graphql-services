package graph

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/vektah/gqlgen/handler"
)

// RegisterRoutes takes a GraphQLServer pointer and registers all needed middlewares and router handlers
func RegisterRoutes(s *GraphQLServer) {
	gqlConfig := Config{
		Resolvers: s,
	}

	// register global middleware to attach a http.ResponseWriter to the context of the request
	s.mux.Use(s.attachResponseWriterMiddleware())

	// register global middleware to attach a user based on the session cookie to the context of the request
	s.mux.Use(s.attachUserMiddleware())

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
			ctx := context.WithValue(r.Context(), CONTEXT_WRITER_KEY, w)

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
			c, err := r.Cookie(SESSION_COOKIE_NAME)

			// allow unauthenticated users in
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			// TODO: validate cookie
			// userId, err := validateAndGetUserID(c)
			// if err != nil {
			// 	http.Error(w, "Invalid cookie", http.StatusForbidden)
			// 	return
			// }

			// get session ID from cookie and attempt to look up the session in redis
			sID := c.Value
			ses, err := s.redisRepository.GetSession(sID)
			if err != nil {
				fmt.Println(err)
			}

			// attach session to context, it might be an empty session if there is an error
			ctx := context.WithValue(r.Context(), CONTEXT_SESSION_KEY, ses)

			// call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// writeSessionCookie takes context and a session ID and users http.ResponseWriter attached to conext to create a cookie
func (s *GraphQLServer) writeSessionCookie(ctx context.Context, sID string) {
	// pull http.Response writer off of context
	w, _ := ctx.Value(CONTEXT_WRITER_KEY).(http.ResponseWriter)

	// create cookie
	c := http.Cookie{
		Name:   SESSION_COOKIE_NAME,
		Value:  sID,
		Domain: s.cfg.Domain,
		// cookie will get expired after 7 days
		Expires: time.Now().AddDate(0, 0, 7),
	}

	// write the cookie to response
	http.SetCookie(w, &c)
}
