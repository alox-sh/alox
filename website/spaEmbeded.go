package website

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"

	"alox.sh"
)

type EmbededSPA struct {
	website
	RootFS    embed.FS
	InputHTML []byte
}

func NewEmbededSPA(inputHTML []byte, rootFS embed.FS) alox.Website {
	return &EmbededSPA{
		website:   newWebsite(),
		RootFS:    rootFS,
		InputHTML: inputHTML,
	}
}

func (embededSPA *EmbededSPA) GetRootFS() fs.FS {
	return embededSPA.RootFS
}

func (embededSPA *EmbededSPA) NewPage(request *http.Request) (*alox.Page, error) {
	return alox.NewPage(embededSPA, request, bytes.NewReader(embededSPA.InputHTML))
}
