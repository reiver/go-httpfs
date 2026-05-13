package httpfs

import (
	"errors"
	"net/http"
	"os"

	"github.com/reiver/go-http204"
	"github.com/reiver/go-http401"
	"github.com/reiver/go-http403"
	"github.com/reiver/go-http404"
	"github.com/reiver/go-http500"
)

func serveHTTPDelete(responseWriter http.ResponseWriter, request *http.Request, root string, authorizerFunc AuthorizerFunc) {
	if nil == responseWriter {
		return
	}
	if nil == request {
		http500.InternalServerError(responseWriter, request)
		return
	}
	if nil == request.URL {
		http500.InternalServerError(responseWriter, request)
		return
	}

	path, ok := PathJoin(root, request.URL.Path)
	if !ok {
		http500.InternalServerError(responseWriter, request)
		return
	}
	if "" == path {
		http500.InternalServerError(responseWriter, request)
		return
	}

	path = pathResolve(root, path)
	if "" == path {
		http404.NotFound(responseWriter, request)
		return
	}

	headers := request.Header
	if nil == headers {
		http500.InternalServerError(responseWriter, request)
		return
	}

	if nil == authorizerFunc {
		http401.Unauthorized(responseWriter, request)
		return
	}
	{
		authorized, err := authorizerFunc(request.Context(), request)
		if nil != err {
			http500.InternalServerError(responseWriter, request)
			return
		}
		if !authorized {
			http403.Forbidden(responseWriter, request)
			return
		}
	}

	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		http404.NotFound(responseWriter, request)
		return
	}
	if nil != err {
		http500.InternalServerError(responseWriter, request)
		return
	}

	http204.NoContent(responseWriter, request)
}
