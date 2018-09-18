package graph

// the key type is needed to avoid warning: "should not use basic type string as key in context.withValue"
type key string

// CONTEXT_SESSION_ID_KEY is a context key used to fetch the session ID from the current request context
const CONTEXT_SESSION_ID_KEY key = "SESSION_ID"

// CONTEXT_SESSION_KEY is a context key used to fetch the session from the current request context
const CONTEXT_SESSION_KEY key = "authed_user"

// CONTEXT_WRITER_KEY is a context key uesd to fetch the http.ResponseWriter from the current request context
const CONTEXT_WRITER_KEY key = "request_writer"

// SESSION_COOKIE_NAME is the name of the session cookie
const SESSION_COOKIE_NAME = "SESSION_ID"
