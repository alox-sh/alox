package main

import (
	"net/http"

	"alox.sh"
	"alox.sh/server"
)

func serveWeb(webServer *server.Web, responseWriter http.ResponseWriter, request *http.Request) {
	_, request.URL.Path = alox.ShiftHead(request.URL.Path)

	if alox.AssertHead(request.URL.Path, func(head string) bool {
		return head == "" || head == "index.html"
	}) {
		// Render HTML page
		return
	}

	// Serve assets
	return
}
