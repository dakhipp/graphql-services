package graph

import (
	"context"
	"net/http"
	"time"
)

// Session is a struct used to identify the user currently making a request
type Session struct {
	ID   string
	Name string
}

// sessionFromUser takes a user and returns a session
func (s *GraphQLServer) sessionFromUser(u *User) Session {
	return Session{
		u.ID,
		u.FirstName + " " + u.LastName,
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
