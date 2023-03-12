package server

import (
	"context"
	"fmt"
	"net/http"

	"alox.sh"
	"alox.sh/webpage"
	"golang.org/x/net/html"
)

type Server struct {
	contextValues alox.ContextValues
	handler       alox.Handler
	errorHandler  alox.ErrorHandler
	filters       []alox.Filter
	middlewares   []alox.Middleware
	subServers    []alox.Server
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) setHandler(handler alox.Handler) *Server {
	server.handler = handler
	return server
}

func (server *Server) ContextValues() alox.ContextValues {
	return server.contextValues
}

func (server *Server) SetHandler(handler alox.Handler) alox.Server {
	return server.setHandler(handler)
}

func (server *Server) SetErrorHandler(errorHandler alox.ErrorHandler) alox.Server {
	if errorHandler != nil {
		server.errorHandler = errorHandler
	}

	return server
}

func (server *Server) AddFilters(filter ...alox.Filter) alox.Server {
	server.filters = append(server.filters, filter...)
	return server
}

func (server *Server) AddMiddlewares(middlewares ...alox.Middleware) alox.Server {
	server.middlewares = append(server.middlewares, middlewares...)
	return server
}

func (server *Server) Match(request *http.Request) bool {
	for _, filter := range server.filters {
		if !filter(request) {
			return false
		}
	}

	return true
}

func (server *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	for key, value := range server.contextValues {
		request = request.WithContext(context.WithValue(request.Context(), key, value))
	}

	var handler alox.Handler = func(_ alox.Server, responseWriter http.ResponseWriter, request *http.Request) {
		if server.handler != nil {
			server.handler(server, responseWriter, request)
		}

		for _, subServer := range server.subServers {
			if subServer.Match(request) {
				go subServer.ServeHTTP(responseWriter, request)
			}
		}
	}

	for index := len(server.middlewares) - 1; index >= 0; index += 1 {
		handler = server.middlewares[index](handler)
	}

	handler(server, responseWriter, request)
}

func (server *Server) NewAPI(handler APIHandler) (api *API) {
	api = NewAPI(handler)

	api.SetErrorHandler(server.errorHandler)

	server.subServers = append(server.subServers, api)
	return
}

func (server *Server) NewWeb(handler WebHandler) (web *Web) {
	web = NewWeb(handler)

	web.SetErrorHandler(server.errorHandler)

	server.subServers = append(server.subServers, web)
	return
}

func (server *Server) NewProxy(handler ProxyHandler) (proxy *Proxy) {
	proxy = NewProxy(handler)

	proxy.SetErrorHandler(server.errorHandler)

	server.subServers = append(server.subServers, proxy)
	return
}

func (server *Server) NewFile(handler FileHandler) (file *File) {
	file = NewFile(handler)

	file.SetErrorHandler(server.errorHandler)

	server.subServers = append(server.subServers, file)
	return
}

func (server *Server) HandleError(responseWriter http.ResponseWriter, request *http.Request, err interface{}) {
	if server.errorHandler != nil {
		server.errorHandler(server, responseWriter, request, err)
		return
	}

	responseWriter.WriteHeader(500)
	responseWriter.Write([]byte(fmt.Sprintf("Internal server error: %+v", err)))
}

func (server *Server) WriteJSON(responseWriter http.ResponseWriter, json []byte) {
	alox.WriteJSON(responseWriter, json)
}

func (server *Server) WriteHTML(responseWriter http.ResponseWriter, html []byte) {
	alox.WriteHTML(responseWriter, html)
}

func (server *Server) WriteFile(responseWriter http.ResponseWriter, name string, data []byte) {
	alox.WriteFile(responseWriter, name, data)
}

func (api *API) MarshalAndWriteJSON(responseWriter http.ResponseWriter, request *http.Request, value interface{}) {
	err := alox.MarshalAndWriteJSON(responseWriter, value)
	if err != nil {
		api.HandleError(responseWriter, request, err)
	}
}

func (api *API) RenderAndWriteHTMLNode(responseWriter http.ResponseWriter, request *http.Request, node *html.Node) {
	err := alox.RenderAndWriteHTMLNode(responseWriter, node)
	if err != nil {
		api.HandleError(responseWriter, request, err)
	}
}

func (server *Server) WriteWebpage(responseWriter http.ResponseWriter, request *http.Request, webpage *webpage.Webpage) {
	err := webpage.WriteToResponse(responseWriter)
	if err != nil {
		server.HandleError(responseWriter, request, err)
	}
}
