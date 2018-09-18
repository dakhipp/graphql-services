package graph

import (
	"context"
	"net/http"
	"time"

	"github.com/dakhipp/graphql-services/auth"
)

// createSessionCookie takes a session ID and returns a cookie with it as the value
func (s *GraphQLServer) createSessionCookie(ctx context.Context, sID string) http.Cookie {
	// cookie will get expired after 7 days
	e := time.Now().AddDate(0, 0, 7)

	// create cookie
	c := http.Cookie{
		Name:    sessionCookieName,
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
		Name:    sessionCookieName,
		Value:   sID,
		Domain:  s.cfg.Domain,
		Expires: time.Now().AddDate(0, 0, -1),
	}

	// return cookie
	return c
}

// writeSessionCookie takes context and a session ID and users http.ResponseWriter attached to conext to create a cookie
func (s *GraphQLServer) writeSessionCookie(ctx context.Context, c http.Cookie) {
	// pull http.Response writer off of context
	w, _ := ctx.Value(contextWriterKey).(http.ResponseWriter)

	// write the cookie to response
	http.SetCookie(w, &c)
}

// getSessionID gets the session ID off of context
func (s *GraphQLServer) getSessionID(ctx context.Context) string {
	// pull session ID off of context
	sID, _ := ctx.Value(contextSessionIDKey).(string)

	return sID
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
