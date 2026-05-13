package httpfs

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/reiver/go-http404"
	"github.com/reiver/go-http405"
	"github.com/reiver/go-http500"
)

func serveHTTPList(responseWriter http.ResponseWriter, request *http.Request, root string) {
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

	fileInfo, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		http404.NotFound(responseWriter, request)
		return
	}
	if nil != err {
		http500.InternalServerError(responseWriter, request)
		return
	}
	if nil == fileInfo {
		http500.InternalServerError(responseWriter, request)
		return
	}
	if !fileInfo.IsDir() {
		http405.MethodNotAllowed(responseWriter, request, methodList)
		return
	}

	entries, err := os.ReadDir(path)
	if nil != err {
		http500.InternalServerError(responseWriter, request)
		return
	}

	requestPath := request.URL.Path
	if !strings.HasSuffix(requestPath, "/") {
		requestPath += "/"
	}

	responseWriter.Header().Set("Content-Type", "text/uri-list")
	responseWriter.WriteHeader(http.StatusOK)

	for _, entry := range entries {
		var buffer [256]byte
		var line []byte = buffer[0:0]

		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		line = append(line, requestPath...)
		line = append(line, name...)

		if entry.IsDir() {
			line = append(line, '/')
		}

		line = append(line, "\r\n"...)

		responseWriter.Write(line)
	}
}
