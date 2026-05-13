package httpfs

import (
	"errors"
	"net/http"
	"os"

	"codeberg.org/reiver/go-ext2media"
	"github.com/reiver/go-http404"
	"github.com/reiver/go-http405"
	"github.com/reiver/go-http500"
	libpath "github.com/reiver/go-path"
)

func serveHTTPGet(responseWriter http.ResponseWriter, request *http.Request, root string) {
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
	if fileInfo.IsDir() {
		http405.MethodNotAllowed(responseWriter, request, methodList)
		return
	}

	file, err := os.Open(path)
	if nil != err {
		http500.InternalServerError(responseWriter, request)
		return
	}
	defer file.Close()

	{
		var mediaType string = "application/octet-stream"

		mt, found := ext2media.Get(libpath.Ext(path))
		if found {
			mediaType = mt
		}

		responseWriter.Header().Set("Content-Type", mediaType)
	}

	http.ServeContent(responseWriter, request, path, fileInfo.ModTime(), file)
}
