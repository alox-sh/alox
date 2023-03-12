package alox

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type Page struct {
	website Website
	request *http.Request

	node *html.Node
	head *html.Node

	shouldRedirect bool
	redirectURI    string

	Title string
	Meta  [][]html.Attribute
}

func NewPage(website Website, request *http.Request, inputHTML io.Reader) (page *Page, err error) {
	page = &Page{
		website: website,
		request: request,
	}

	if page.node, err = html.Parse(inputHTML); err != nil {
		return
	}

	return
}

func (page *Page) MarshalMap() map[string]interface{} {
	return map[string]interface{}{
		"title": page.Title,
		"meta":  page.Meta,
		// "data":  page.Data,
	}
}

func (page *Page) WriteToResponse(responseWriter http.ResponseWriter) (err error) {
	if err = page.enrichHTML(); err != nil {
		return
	}

	responseHeader := responseWriter.Header()

	if page.shouldRedirect && page.website.GetDisableHTTPRedirect() {
		responseHeader.Set("Location", page.redirectURI)
		responseWriter.WriteHeader(http.StatusFound)
		return
	}

	responseHeader.Set("Content-Type", "text/html; charset=utf-8")
	responseWriter.WriteHeader(http.StatusOK)

	return html.Render(responseWriter, page.node)
}

func (page *Page) WriteToBuffer(output *bytes.Buffer) (err error) {
	if err = page.enrichHTML(); err != nil {
		return
	}

	return html.Render(output, page.node)
}

func (page *Page) WriteToBytes() (output []byte, err error) {
	var buffer bytes.Buffer

	if err = page.WriteToBuffer(&buffer); err != nil {
		return
	}

	return buffer.Bytes(), nil
}

func (page *Page) enrichHTML() (err error) {
	if page.head = page.findHeadNode(); page.head == nil {
		return fmt.Errorf("Invalid HTML input: missing head element")
	}

	if err = page.injectState(); err != nil {
		return
	}

	if page.shouldRedirect {
		page.setRedirect()
	}

	if len(page.Title) > 0 {
		page.setTitle()
	}

	if len(page.Meta) > 0 {
		page.setMeta()
	}

	return
}

func (page *Page) injectState() (err error) {
	var state string

	if state, err = page.website.MarshalStateJS(page); err != nil {
		return
	}

	scriptNode := &html.Node{
		Type: html.ElementNode,
		Data: "script",
		Attr: []html.Attribute{},
	}

	scriptNode.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: state,
	})

	page.head.AppendChild(scriptNode)
	return
}

func (page *Page) setRedirect() {
	scriptNode := &html.Node{
		Type: html.ElementNode,
		Data: "script",
		Attr: []html.Attribute{},
	}

	scriptNode.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: fmt.Sprintf(
			"if (window && window.location && typeof window.location.assign === 'function') {\n    window.location.assign(\"%s\");\n}\n",
			page.redirectURI,
		),
	})

	page.head.InsertBefore(scriptNode, page.head.FirstChild)
}

func (page *Page) setTitle() {
	titleNode := page.findTitleNode()

	if len(page.Title) > 0 {
		titleNodeText := &html.Node{
			Type: html.TextNode,
			Data: page.Title,
		}

		if titleNode == nil {
			titleNode := &html.Node{
				Type: html.ElementNode,
				Data: "title",
				Attr: []html.Attribute{},
			}

			titleNode.AppendChild(titleNodeText)
			page.head.InsertBefore(titleNode, page.head.FirstChild)
		} else {
			for child := titleNode.FirstChild; child != nil; child = child.NextSibling {
				titleNode.RemoveChild(child)
			}

			titleNode.AppendChild(titleNodeText)
		}
	}
}

func (page *Page) setMeta() {
	for _, meta := range page.Meta {
		page.head.InsertBefore(&html.Node{
			Type: html.ElementNode,
			Data: "meta",
			Attr: meta,
		}, page.head.FirstChild)
	}
}

func (page *Page) findHeadNode() *html.Node {
	return findElement(page.node, "head")
}

func (page *Page) findTitleNode() *html.Node {
	return findElement(page.node, "title")
}
