package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	*Proxy
	*httputil.ReverseProxy
	OriginURL *url.URL
}

func NewReverseProxy(originURL *url.URL) (reverseProxy *ReverseProxy, err error) {
	if originURL == nil {
		originURL = &url.URL{}
	}

	reverseProxy = &ReverseProxy{
		OriginURL:    originURL,
		ReverseProxy: httputil.NewSingleHostReverseProxy(originURL),
		Proxy: NewProxy(func(_ *Proxy, responseWriter http.ResponseWriter, request *http.Request) {
			reverseProxy.ServeHTTP(responseWriter, request)
		}),
	}
	return
}

func (reverseProxy *ReverseProxy) SetOriginURL(originURL *url.URL) {
	reverseProxy.OriginURL = originURL
	reverseProxy.ReverseProxy = httputil.NewSingleHostReverseProxy(originURL)
}
