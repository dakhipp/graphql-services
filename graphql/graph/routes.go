package graph

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dakhipp/graphql-services/auth"
	"github.com/vektah/gqlgen/handler"
)

// RegisterRoutes takes a GraphQLServer pointer and registers all needed middlewares, router handlers, and directives
func RegisterRoutes(s *GraphQLServer) {
	gqlConfig := Config{
		Resolvers: s,
	}

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

			// get session ID and attempt to look up the session in redis
			sID := c.Value
			ses, err := s.redisRepository.GetSession(sID)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(fmt.Sprintf("Currently logged in user: %v", ses))

			// attach session to context, it might be an empty session if there is an error
			ctx := context.WithValue(r.Context(), CONTEXT_SESSION_KEY, ses)

			// attach session ID to context, this is needed in order to delete the session
			ctx2 := context.WithValue(ctx, CONTEXT_SESSION_ID_KEY, sID)

			// call the next with our new context
			r = r.WithContext(ctx)
			r = r.WithContext(ctx2)
			next.ServeHTTP(w, r)
		})
	}
}

// createSessionCookie takes a session ID and returns a cookie with it as the value
func (s *GraphQLServer) createSessionCookie(ctx context.Context, sID string) http.Cookie {
	// cookie will get expired after 7 days
	e := time.Now().AddDate(0, 0, 7)

	// create cookie
	c := http.Cookie{
		Name:    SESSION_COOKIE_NAME,
		Value:   sID,
		Domain:  s.cfg.Domain,
		Expires: e,
	}

	// return cookie
	return c
}

// expireSessionCookie pulls the session ID off of context and creates a new cookie expired cookie using it as the value
func (s *GraphQLServer) expireSessionCookie(sID string) http.Cookie {
	// create cookie
	c := http.Cookie{
		Name:    SESSION_COOKIE_NAME,
		Value:   sID,
		Domain:  s.cfg.Domain,
		Expires: time.Now().AddDate(0, 0, -1),
	}

	// return cookie
	return c
}

// getSessionID gets the session ID off of context
func (s *GraphQLServer) getSessionID(ctx context.Context) string {
	// pull session ID off of context
	sID, _ := ctx.Value(CONTEXT_SESSION_ID_KEY).(string)

	return sID
}

// writeSessionCookie takes context and a session ID and users http.ResponseWriter attached to conext to create a cookie
func (s *GraphQLServer) writeSessionCookie(ctx context.Context, c http.Cookie) {
	// pull http.Response writer off of context
	w, _ := ctx.Value(CONTEXT_WRITER_KEY).(http.ResponseWriter)

	// write the cookie to response
	http.SetCookie(w, &c)
}

// createSession takes an authenticated user response and returns a session as well as a session cookie
func (s *GraphQLServer) createSession(ctx context.Context, resp *auth.User) *Session {
	// create unique session ID and a session based on the user who authenticated
	ses := &Session{
		ID:    resp.ID,
		Roles: toRoles(resp.Roles),
	}

	// return the session
	return ses
}

// toRoles converts a string array into a Roles array
func toRoles(s []string) []Role {
	c := make([]Role, len(s))
	for i, v := range s {
		c[i] = Role(v)
	}
	return c
}
