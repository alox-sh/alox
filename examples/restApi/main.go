package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"alox.sh"
	"alox.sh/server"

	"examples/restApi/config"
)

const HTTPPort = 2000

func main() {
	var addr = fmt.Sprintf(":%d", HTTPPort)

	config, err := config.SetupConfig()
	if err != nil {
		log.Fatal(err)
	}

	appServer := server.NewServer()

	appServer.SetContext(config.ToContext(context.Background()))

	appServer.
		NewAPI(serveAPI).
		AddFilters(func(request *http.Request) bool {
			return alox.MatchHead(request.URL.Path, "api")
		})

	appServer.
		NewWeb(serveWeb).
		AddFilters(func(request *http.Request) bool {
			return alox.MatchHead(request.URL.Path, "web")
		})

	fmt.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, appServer)
}
