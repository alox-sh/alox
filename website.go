package alox

import (
	"io/fs"
	"net/http"
)

type Website interface {
	GetRootFS() fs.FS
	GetDisableHTTPRedirect() bool
	GetStateGlobalName() string

	SetDisableHTTPRedirect(bool)
	SetStateGlobalName(string)

	NewPage(*http.Request) (*Page, error)

	MarshalStateMap(*Page) map[string]interface{}
	MarshalStateJS(*Page) (string, error)

	ServeHTTP(http.ResponseWriter, *http.Request)
}
