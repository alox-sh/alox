package filter

import (
	"net/http"
	"strings"

	"alox.sh"
)

func Methods(methods ...string) alox.Filter {
	return func(request *http.Request) bool {
		methodLower := strings.ToLower(request.Method)

		for _, method := range methods {
			if methodLower == strings.ToLower(method) {
				return true
			}
		}

		return false
	}
}

var (
	GET     = Methods("get")
	POST    = Methods("post")
	PUT     = Methods("put")
	PATCH   = Methods("patch")
	DELETE  = Methods("delete")
	OPTIONS = Methods("options")
)
