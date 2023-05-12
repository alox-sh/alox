package server

import (
	"alox.sh"
)

type Router struct {
	*Server
}

func NewRouter() (api *API) {
	api = &API{Server: NewServer()}
	api.Server.setHandler(nil)
	return
}

func (router *Router) NewRoute(route alox.Server, filters ...alox.Filter) {
	route.AddFilters(filters...)
	router.Sub(route)
}
