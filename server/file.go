package server

import (
	"net/http"

	"alox.sh"
)

type FileHandler func(*File, http.ResponseWriter, *http.Request)

type File struct {
	*Server
}

func NewFile(handler FileHandler) (file *File) {
	file = &File{Server: NewServer()}
	return file.setHandler(handler)
}

func (file *File) setHandler(handler FileHandler) *File {
	file.Server.setHandler(func(_ alox.Server, responseWriter http.ResponseWriter, request *http.Request) {
		handler(file, responseWriter, request)
	})
	return file
}
