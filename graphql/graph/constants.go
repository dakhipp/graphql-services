package graph

// the key type is needed to avoid warning: "should not use basic type string as key in context.withValue"
type key string

// contextSessionIDKey is a context key used to fetch the session ID from the current request context
const contextSessionIDKey key = "SESSION_ID"

// contextSessionKey is a context key used to fetch the session from the current request context
const contextSessionKey key = "authed_user"

// contextWriterKey is a context key uesd to fetch the http.ResponseWriter from the current request context
const contextWriterKey key = "request_writer"

// sessionCookieName is the name of the session cookie
const sessionCookieName = "SESSION_ID"
