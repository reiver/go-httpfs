package httpfs

import (
	"net/http"
	"strings"
	"time"

	"github.com/reiver/go-http405"
	"github.com/reiver/go-http500"
)

type ReadWriterHandler struct {
	RootDir               string
	AuthorizerFunc        AuthorizerFunc
	HTTPBodyReadSizeLimit int64
	HTTPBodyReadTimeOut   time.Duration
}

func (receiver ReadWriterHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if nil == responseWriter {
		return
	}
	if nil == request {
		http500.InternalServerError(responseWriter, request)
		return
	}

	switch strings.ToUpper(request.Method) {
	case methodCreate:
		serveHTTPCreate(responseWriter, request, receiver.RootDir, receiver.AuthorizerFunc, receiver.HTTPBodyReadSizeLimit, receiver.HTTPBodyReadTimeOut)
		return
	case http.MethodDelete:
		serveHTTPDelete(responseWriter, request, receiver.RootDir, receiver.AuthorizerFunc)
		return
	case http.MethodGet:
		serveHTTPGet(responseWriter, request, receiver.RootDir)
		return
	case http.MethodHead:
		serveHTTPHead(responseWriter, request, receiver.RootDir)
		return
	case methodList:
		serveHTTPList(responseWriter, request, receiver.RootDir)
		return
	case http.MethodPut:
		serveHTTPPut(responseWriter, request, receiver.RootDir, receiver.AuthorizerFunc, receiver.HTTPBodyReadSizeLimit, receiver.HTTPBodyReadTimeOut)
		return
	default:
		http405.MethodNotAllowed(responseWriter, request,
			methodCreate,
			http.MethodDelete,
			http.MethodGet,
			http.MethodHead,
			methodList,
			http.MethodPut,
		)
		return
	}
}
