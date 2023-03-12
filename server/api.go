package server

import (
	"net/http"

	"alox.sh"
)

type APIHandler func(*API, http.ResponseWriter, *http.Request)

type API struct {
	*Server
}

func NewAPI(handler APIHandler) (api *API) {
	api = &API{Server: NewServer()}
	return api.setHandler(handler)
}

func (api *API) setHandler(handler APIHandler) *API {
	api.Server.setHandler(func(_ alox.Server, responseWriter http.ResponseWriter, request *http.Request) {
		handler(api, responseWriter, request)
	})
	return api
}
