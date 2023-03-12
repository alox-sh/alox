package server

import (
	"net/http"

	"alox.sh"
)

type ProxyHandler func(*Proxy, http.ResponseWriter, *http.Request)

type Proxy struct {
	*Server
}

func NewProxy(handler ProxyHandler) (proxy *Proxy) {
	proxy = &Proxy{Server: NewServer()}
	return proxy.setHandler(handler)
}

func (proxy *Proxy) setHandler(handler ProxyHandler) *Proxy {
	proxy.Server.setHandler(func(_ alox.Server, responseWriter http.ResponseWriter, request *http.Request) {
		handler(proxy, responseWriter, request)
	})
	return proxy
}
