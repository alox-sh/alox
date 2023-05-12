package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html"

	"alox.sh"
	"alox.sh/webpage"
)

type contextValues map[interface{}]interface{}

func (contextValues contextValues) inject(request *http.Request) *http.Request {
	requestContext := request.Context()

	for key, value := range contextValues {
		requestContext = context.WithValue(requestContext, key, value)
	}

	return request.WithContext(requestContext)
}

func (contextValues contextValues) Get(key interface{}) interface{} {
	return contextValues[key]
}

func (contextValues contextValues) Set(key, value interface{}) {
	contextValues[key] = value
}

func (contextValues contextValues) Del(key interface{}) {
	delete(contextValues, key)
}

type Server struct {
	handler       alox.Handler
	errorHandler  alox.ErrorHandler
	onWrite       alox.OnWrite
	contextValues contextValues
	filters       []alox.Filter
	middlewares   []alox.Middleware
	sub           []alox.Server
}

func NewServer() *Server {
	return &Server{contextValues: contextValues{}}
}

func (server *Server) setHandler(handler alox.Handler) *Server {
	server.handler = handler
	return server
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

func (server *Server) ContextValues() alox.ContextValues {
	return server.contextValues
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
	request = server.contextValues.inject(request)

	var handler alox.Handler = func(_ alox.Server, responseWriter http.ResponseWriter, request *http.Request) {
		if server.handler != nil {
			go server.handler(server, responseWriter, request)
		}

		for _, sub := range server.sub {
			if sub.Match(request) {
				sub.ServeHTTP(responseWriter, request)
			}
		}
	}

	for index := len(server.middlewares) - 1; index >= 0; index -= 1 {
		handler = server.middlewares[index](handler)
	}

	handler(server, responseWriter, request)
}

func (server *Server) NewAPI(handler APIHandler) (api *API) {
	api = NewAPI(handler)

	api.SetErrorHandler(server.errorHandler)

	server.sub = append(server.sub, api)
	return
}

func (server *Server) NewWeb(handler WebHandler) (web *Web) {
	web = NewWeb(handler)

	web.SetErrorHandler(server.errorHandler)

	server.sub = append(server.sub, web)
	return
}

func (server *Server) NewSPA(handler SPAHandler, params SPAParams) (spa *SPA, err error) {
	if spa, err = NewSPA(handler, params); err != nil {
		return
	}

	spa.SetErrorHandler(server.errorHandler)

	server.sub = append(server.sub, spa)
	return
}

// func (server *Server) NewProxy(handler ProxyHandler) (proxy *Proxy) {
// 	proxy = NewProxy(handler)

// 	proxy.SetErrorHandler(server.errorHandler)

// 	server.subServers = append(server.subServers, proxy)
// 	return
// }

func (server *Server) NewReverseProxy(originURL *url.URL) (reverseProxy *ReverseProxy, err error) {
	if reverseProxy, err = NewReverseProxy(originURL); err != nil {
		return
	}

	reverseProxy.SetErrorHandler(server.errorHandler)

	server.sub = append(server.sub, reverseProxy)
	return
}

func (server *Server) NewFile(handler FileHandler) (file *File) {
	file = NewFile(handler)

	file.SetErrorHandler(server.errorHandler)

	server.sub = append(server.sub, file)
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

func (server *Server) OnWrite(onWrite alox.OnWrite) {
	server.onWrite = onWrite
}

func (server *Server) Write(
	responseWriter http.ResponseWriter,
	request *http.Request,
	contentType string,
	contentLength int,
	data []byte,
) {
	if server.onWrite != nil {
		go server.onWrite(request, contentType, contentLength, &data)
	}

	alox.Write(responseWriter, contentType, contentLength, data)
}

func (server *Server) WriteText(responseWriter http.ResponseWriter, request *http.Request, text []byte) {
	server.Write(responseWriter, request, "text/plain", len(text), text)
}

func (server *Server) WriteHTML(responseWriter http.ResponseWriter, request *http.Request, html []byte) {
	server.Write(responseWriter, request, "text/html", len(html), html)
}

func (server *Server) WriteJSON(responseWriter http.ResponseWriter, request *http.Request, json []byte) {
	server.Write(responseWriter, request, "application/json", len(json), json)
}

func (server *Server) WriteXML(responseWriter http.ResponseWriter, request *http.Request, xml []byte) {
	server.Write(responseWriter, request, "application/xml", len(xml), xml)
}

// func (server *Server) WriteFile(responseWriter http.ResponseWriter, name string, data []byte) {
// 	alox.WriteFile(responseWriter, name, data)
// }

func (server *Server) MarshalAndWriteJSON(responseWriter http.ResponseWriter, request *http.Request, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		server.HandleError(responseWriter, request, err)
		return
	}

	server.WriteJSON(responseWriter, request, data)
}

func (server *Server) RenderAndWriteHTMLNode(responseWriter http.ResponseWriter, request *http.Request, node *html.Node) {
	buffer := &bytes.Buffer{}

	if err := html.Render(buffer, node); err != nil {
		server.HandleError(responseWriter, request, err)
		return
	}

	server.WriteHTML(responseWriter, request, buffer.Bytes())
}

func (server *Server) WriteWebpage(responseWriter http.ResponseWriter, request *http.Request, webpage *webpage.Webpage) (err error) {
	buffer := &bytes.Buffer{}

	if err = webpage.WriteToBuffer(buffer); err != nil {
		return
	}

	contentType := "text/html"
	if len(webpage.Charset) > 0 {
		contentType = fmt.Sprintf("%s; charset=%s", contentType, webpage.Charset)
	}

	responseWriter.WriteHeader(http.StatusOK)

	if request.Method != "HEAD" {
		server.Write(responseWriter, request, contentType, buffer.Len(), buffer.Bytes())
	}

	return
}

func (server *Server) MustWriteWebpage(responseWriter http.ResponseWriter, request *http.Request, webpage *webpage.Webpage) {
	err := server.WriteWebpage(responseWriter, request, webpage)
	if err != nil {
		server.HandleError(responseWriter, request, err)
	}
}
