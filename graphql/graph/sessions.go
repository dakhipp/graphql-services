package graph

import (
	"context"
	"net/http"
	"time"

	"github.com/dakhipp/graphql-services/auth/pb"
	"github.com/segmentio/ksuid"
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
func (s *GraphQLServer) createSession(ctx context.Context, resp *pb.AuthResponse) *Session {
	// create unique session ID and a session based on the user who authenticated
	ses := &Session{
		ID:            resp.User.Id,
		FirstName:     resp.User.FirstName,
		LastName:      resp.User.LastName,
		Email:         resp.User.Email,
		Phone:         resp.User.Phone,
		Roles:         toRoles(resp.User.Roles),
		EmailVerified: resp.User.EmailVerified,
		PhoneVerified: resp.User.PhoneVerified,
	}

	// return the session
	return ses
}

// handleWriteSession creates a session ID, creates a session, creates and writes a session cookie, and saves the session in redis. It is called at the end of the login and register mutations.
func (s *GraphQLServer) handleWriteSession(ctx context.Context, resp *pb.AuthResponse) (*Session, error) {
	// create unique session ID
	sID := ksuid.New().String()

	// create session
	ses := s.createSession(ctx, resp)

	// create session cookie
	c := s.createSessionCookie(ctx, sID)

	// use http.ResponseWriter to write the cookie into the response
	s.writeSessionCookie(ctx, c)

	// save session in redis
	s.redisRepository.CreateSession(sID, ses)

	// return the session
	return ses, nil
}

// getSessionFromContext takes context returns the current session
func (s *GraphQLServer) getSessionFromContext(ctx context.Context) Session {
	// pull user off of context and return the user
	ses, _ := ctx.Value(contextSessionKey).(Session)
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
