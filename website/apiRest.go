package website

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"

	"alox.sh"
)

type RestAPI struct {
	website
}

func NewRestAPI() alox.Website {
	return &RestAPI{
		website: newWebsite(),
	}
}

func (restAPI *RestAPI) GetRootFS() fs.FS {
	return os.DirFS(".")
}

func (restAPI *RestAPI) NewPage(request *http.Request) (*alox.Page, error) {
	return alox.NewPage(restAPI, request, bytes.NewReader([]byte("")))
}
