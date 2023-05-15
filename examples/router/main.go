package main

import (
	"fmt"
	"net/http"

	"alox.sh"
	"alox.sh/server"
)

const addr = ":8080"

var API = server.NewAPI(func(api *server.API, responseWriter http.ResponseWriter, request *http.Request) {
	api.MarshalAndWriteJSON(responseWriter, request, map[string]interface{}{
		"message": "Hello from the API!",
	})
})

var Web = server.NewWeb(func(web *server.Web, responseWriter http.ResponseWriter, request *http.Request) {
	web.WriteHTML(responseWriter, request, []byte(`
<!DOCTYPE html>
<html>
	<body>
		Hello from the web!
	</body>
</html>
`))
})

func main() {
	router := server.NewRouter()

	router.NewRoute(API, func(request *http.Request) bool {
		return alox.MatchHead(request.URL.Path, "api")
	})

	router.NewRoute(Web, func(request *http.Request) bool {
		return alox.MatchHead(request.URL.Path, "web")
	})

	router.NewRoute(server.NewAPI(func(api *server.API, responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusNotFound)
		api.WriteJSON(responseWriter, request, []byte(`{"error":"Not Found"}`))
	}), func(request *http.Request) bool {
		return !alox.MatchHead(request.URL.Path, "api") && !alox.MatchHead(request.URL.Path, "web")
	})

	fmt.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, router)
}
