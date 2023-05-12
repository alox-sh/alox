package alox

import (
	"net/http"
)

type (
	Filter       func(*http.Request) bool
	Handler      func(Server, http.ResponseWriter, *http.Request)
	ErrorHandler func(Server, http.ResponseWriter, *http.Request, interface{})
	Middleware   func(Handler) Handler
	OnWrite      func(*http.Request, string, int, *[]byte)
)

type ContextValues interface {
	Get(key interface{}) interface{}
	Set(key, value interface{})
	Del(key interface{})
}

type ResponseMethods interface {
	Write(
		responseWriter http.ResponseWriter,
		request *http.Request,
		contentType string,
		contentLength int,
		data []byte,
	)

	WriteText(http.ResponseWriter, *http.Request, []byte)
	WriteHTML(http.ResponseWriter, *http.Request, []byte)
	WriteJSON(http.ResponseWriter, *http.Request, []byte)
	WriteXML(http.ResponseWriter, *http.Request, []byte)

	OnWrite(OnWrite)
}

type Server interface {
	http.Handler
	ResponseMethods

	SetHandler(Handler) Server
	SetErrorHandler(ErrorHandler) Server

	ContextValues() ContextValues

	AddFilters(...Filter) Server
	AddMiddlewares(...Middleware) Server

	Match(*http.Request) bool

	HandleError(http.ResponseWriter, *http.Request, interface{})
}
