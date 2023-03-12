package server

import (
	"net/http"
)

type ForwardProxy struct {
	*Proxy
}

func NewForwardProxy() (forwardProxy *ForwardProxy) {
	forwardProxy = &ForwardProxy{
		Proxy: NewProxy(func(_ *Proxy, responseWriter http.ResponseWriter, request *http.Request) {
		}),
	}

	return
}
