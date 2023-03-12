package main

import (
	"net/http"

	"alox.sh"
	"alox.sh/filter"
	"alox.sh/server"

	"examples/restApi/models"
)

var endpoints = struct {
	getTodos  server.EndpointHandler
	postTodos server.EndpointHandler
	getTodo   server.EndpointHandler
}{
	getTodos: func(endpoint *server.Endpoint, responseWriter http.ResponseWriter, request *http.Request) {
		request.Context()

		models.GetTodos()
	},
	postTodos: func(endpoint *server.Endpoint, responseWriter http.ResponseWriter, request *http.Request) {},
	getTodo:   func(endpoint *server.Endpoint, responseWriter http.ResponseWriter, request *http.Request) {},
}

func serveAPI(apiServer *server.API, responseWriter http.ResponseWriter, request *http.Request) {
	var head string

	// remove the /api prefix
	_, request.URL.Path = alox.ShiftHead(request.URL.Path)
	head, request.URL.Path = alox.ShiftHead(request.URL.Path)

	switch head {
	case "todos":
		head, request.URL.Path = alox.ShiftHead(request.URL.Path)

		apiServer.
			NewEndpoint(endpoints.getTodos).
			AddFilters(filter.GET, filter.PathSegments().LenEq(0))

		apiServer.
			NewEndpoint(endpoints.postTodos).
			AddFilters(filter.POST, filter.PathSegments().LenEq(0))

		apiServer.
			NewEndpoint(endpoints.getTodo).
			AddFilters(filter.GET, filter.PathSegments().LenEq(1))

		return
	case "tasks":
		return
	}
}
