package server

import (
	"alox.sh"
)

type Router struct {
	*Server
}

func NewRouter() (router *Router) {
	router = &Router{Server: NewServer()}
	router.Server.setHandler(nil)
	return
}

func (router *Router) NewRoute(route alox.Server, filters ...alox.Filter) {
	route.AddFilters(filters...)
	router.Sub(route)
}
