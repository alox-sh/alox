package alox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type OnWrite func(request *http.Request, contentType string, contentLength int, data *[]byte)

type ResponseMethods interface {
	Write(
		responseWriter http.ResponseWriter,
		request *http.Request,
		contentType string,
		contentLength int,
		data []byte,
	)

	WriteText(responseWriter http.ResponseWriter, request *http.Request, text []byte)
	WriteHTML(responseWriter http.ResponseWriter, request *http.Request, html []byte)
	WriteJSON(responseWriter http.ResponseWriter, request *http.Request, json []byte)
	WriteXML(responseWriter http.ResponseWriter, request *http.Request, xml []byte)

	OnWrite(onWrite OnWrite)
}

func Write(
	responseWriter http.ResponseWriter,
	contentType string,
	contentLength int,
	data []byte,
) {
	responseWriter.Header().Set("Content-Type", contentType)
	responseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	responseWriter.Write(data)
}

func WriteText(responseWriter http.ResponseWriter, text []byte) {
	Write(responseWriter, "text/plain", len(text), text)
}

func WriteHTML(responseWriter http.ResponseWriter, html []byte) {
	Write(responseWriter, "text/html", len(html), html)
}

func WriteJSON(responseWriter http.ResponseWriter, json []byte) {
	Write(responseWriter, "application/json", len(json), json)
}

func WriteXML(responseWriter http.ResponseWriter, xml []byte) {
	Write(responseWriter, "application/xml", len(xml), xml)
}

func WriteFile(responseWriter http.ResponseWriter, contentType string, data []byte) {
	header := responseWriter.Header()
	header.Set("Content-Type", contentType)
	header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

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
	buffer := &bytes.Buffer{}

	if err = html.Render(buffer, node); err != nil {
		return
	}

	header := responseWriter.Header()
	header.Set("Content-Type", "text/html")
	header.Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))

	_, err = responseWriter.Write(buffer.Bytes())
	return
}
