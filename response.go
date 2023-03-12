package alox

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type ResponseMethods interface {
	WriteJSON(responseWriter http.ResponseWriter, json []byte)
	WriteHTML(responseWriter http.ResponseWriter, html []byte)
}

func WriteJSON(responseWriter http.ResponseWriter, json []byte) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(json)
}

func WriteHTML(responseWriter http.ResponseWriter, html []byte) {
	responseWriter.Header().Set("Content-Type", "text/html")
	responseWriter.Write(html)
}

func WriteFile(responseWriter http.ResponseWriter, name string, data []byte) {
	header := responseWriter.Header()

	nameParts := strings.Split(name, ".")
	header.Set("Content-Type", mime.TypeByExtension("."+nameParts[len(nameParts)-1]))

	header.Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))

	responseWriter.Write(data)
}

func MarshalAndWriteJSON(responseWriter http.ResponseWriter, value interface{}) (err error) {
	var data []byte

	if data, err = json.Marshal(value); err != nil {
		return
	}

	WriteJSON(responseWriter, data)
	return
}

func RenderAndWriteHTMLNode(responseWriter http.ResponseWriter, node *html.Node) (err error) {
	responseWriter.Header().Set("Content-Type", "text/html")
	return html.Render(responseWriter, node)
}
