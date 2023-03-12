package alox

import (
	"net/http"
)

type Filter func(*http.Request) bool

type Handler func(Server, http.ResponseWriter, *http.Request)

type ErrorHandler func(Server, http.ResponseWriter, *http.Request, interface{})

type Middleware func(Handler) Handler

type ContextValues map[interface{}]interface{}

func (contextValues ContextValues) Get(key interface{}) interface{} {
	return contextValues[key]
}

func (contextValues ContextValues) Set(key, value interface{}) {
	contextValues[key] = value
}

func (contextValues ContextValues) Del(key interface{}) {
	delete(contextValues, key)
}

type Server interface {
	http.Handler
	ResponseMethods

	SetHandler(handler Handler) Server

	ContextValues() ContextValues // map[interface{}]interface{}
	// SetContextValue(key, value interface{})

	// // Context returns the base context.Context.
	// //
	// // This context is propagated to all sub-servers and its
	// Context() context.Context
	// SetContext(context context.Context) Server

	AddFilters(filters ...Filter) Server
	AddMiddlewares(middlewares ...Middleware) Server

	// Match checks whether incoming *http.Request should be handled
	// by this server's handler.
	//
	// Match function will test request against all filters associated
	// with this server. If the request doesn't pass any one of them,
	// the Match function will return false. Otherwise, by default,
	// the Match function will return true.
	Match(request *http.Request) bool
}
