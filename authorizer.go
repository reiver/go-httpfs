package httpfs

import (
	"context"
	"net/http"
)

// AuthorizerFunc is a function that return true if an HTTP request is authorized and returns false if it isn't.
//
// See also:
//
//	• [ReadWriterHandler]
type AuthorizerFunc func(context.Context, *http.Request) (bool, error)
