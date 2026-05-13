package httpfs

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/reiver/go-http204"
	"github.com/reiver/go-http401"
	"github.com/reiver/go-http403"
	"github.com/reiver/go-http404"
	"github.com/reiver/go-http500"
)

func serveHTTPAppend(responseWriter http.ResponseWriter, request *http.Request, root string, authorizerFunc AuthorizerFunc, httpBodyReadSizeLimit int64, httpBodyReadTimeOut time.Duration) {
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

	path = pathResolve(root, path)
	if "" == path {
		http404.NotFound(responseWriter, request)
		return
	}

	{
		if httpBodyReadTimeOut <= 0 {
			httpBodyReadTimeOut = 2*time.Minute
		}
		err := http.NewResponseController(responseWriter).SetReadDeadline(time.Now().Add(httpBodyReadTimeOut))
		if nil != err {
			ctx, cancel := context.WithTimeout(request.Context(), httpBodyReadTimeOut)
			defer cancel()
			request = request.WithContext(ctx)
		}
	}

	if httpBodyReadSizeLimit <= 0 {
		httpBodyReadSizeLimit = 1_073_741_824 // 2^30
	}
	body := http.MaxBytesReader(responseWriter, request.Body, httpBodyReadSizeLimit)
	defer body.Close()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if nil != err {
		http500.InternalServerError(responseWriter, request)
		return
	}

	_, err = io.Copy(file, body)
	if closeErr := file.Close(); nil == err {
		err = closeErr
	}
	if nil != err {
		http500.InternalServerError(responseWriter, request)
		return
	}

	http204.NoContent(responseWriter, request)
}
