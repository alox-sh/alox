package webpage

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

// icons / favicons / etc. ?
//
// https://careerkarma.com/blog/html-favicon
// https://www.30secondsofcode.org/articles/s/html-head-icons

type Webpage struct {
	rootNode *html.Node
	headNode *html.Node
	bodyNode *html.Node

	redirectURI  string
	noscriptNode *html.Node

	Meta        [][]html.Attribute
	Links       [][]html.Attribute
	Styles      []Style
	Scripts     []Script
	Robots      string
	Viewport    string
	Charset     string
	Title       string
	Description string
	Author      string
}

type Style struct {
	Content    string
	Attributes []html.Attribute
}

type Script struct {
	Src        string
	Content    string
	Head       bool
	Attributes []html.Attribute
}

func newWebpage() *Webpage {
	return &Webpage{
		noscriptNode: newNoscriptNode(),
		Charset:      "utf-8",
		Viewport:     "width=device-width, initial-scale=1",
	}
}

func NewWebpage(inputHTML io.Reader) (webpage *Webpage, err error) {
	webpage = newWebpage()

	if webpage.rootNode, err = html.Parse(inputHTML); err != nil {
		return
	}

	err = webpage.populate()
	return
}

func (webpage *Webpage) Clone() (*Webpage, error) {
	clone := newWebpage()
	clone.rootNode = cloneNode(webpage.rootNode)

	clone.Meta = make([][]html.Attribute, len(webpage.Meta))
	copy(clone.Meta, webpage.Meta)

	clone.Links = make([][]html.Attribute, len(webpage.Links))
	copy(clone.Links, webpage.Links)

	clone.Styles = make([]Style, len(webpage.Styles))
	copy(clone.Styles, webpage.Styles)

	clone.Scripts = make([]Script, len(webpage.Scripts))
	copy(clone.Scripts, webpage.Scripts)

	clone.Robots = webpage.Robots
	clone.Viewport = webpage.Viewport
	clone.Charset = webpage.Charset
	clone.Title = webpage.Title
	clone.Description = webpage.Description
	clone.Author = webpage.Author

	return clone, clone.populate()
}

func (webpage *Webpage) SetRedirect(redirectURI string) {
	webpage.redirectURI = redirectURI
}

func (webpage *Webpage) WriteToResponse(responseWriter http.ResponseWriter, request *http.Request) (err error) {
	webpage.enrichHTML()

	header := responseWriter.Header()

	if len(webpage.redirectURI) > 0 {
		header.Set("Location", webpage.redirectURI)
		responseWriter.WriteHeader(http.StatusFound)
		return
	}

	buffer := &bytes.Buffer{}
	if err = html.Render(buffer, webpage.rootNode); err != nil {
		return
	}

	contentType := "text/html"
	if len(webpage.Charset) > 0 {
		contentType = fmt.Sprintf("%s; charset=%s", contentType, webpage.Charset)
	}

	header.Set("Content-Length", fmt.Sprintf("%d", buffer.Len()))
	header.Set("Content-Type", contentType)

	responseWriter.WriteHeader(http.StatusOK)

	if request.Method != "HEAD" {
		_, err = responseWriter.Write(buffer.Bytes())
	}

	return
}

func (webpage *Webpage) WriteToBuffer(outputHTML *bytes.Buffer) (err error) {
	webpage.enrichHTML()

	return html.Render(outputHTML, webpage.rootNode)
}

func (webpage *Webpage) WriteToBytes() (outputHTML []byte, err error) {
	var buffer bytes.Buffer

	if err = webpage.WriteToBuffer(&buffer); err != nil {
		return
	}

	return buffer.Bytes(), nil
}

func (webpage *Webpage) populate() (err error) {
	if webpage.headNode = webpage.findHeadNode(); webpage.headNode == nil {
		err = fmt.Errorf("Invalid HTML input: missing head element")
		return
	}

	if webpage.bodyNode = webpage.findBodyNode(); webpage.bodyNode == nil {
		err = fmt.Errorf("Invalid HTML input: missing body element")
		return
	}

processHead:
	for child := webpage.headNode.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "title" {
			for subChild := child.FirstChild; subChild != nil; subChild = subChild.NextSibling {
				if subChild.Type == html.TextNode {
					webpage.Title = subChild.Data
					webpage.headNode.RemoveChild(child)
					continue processHead
				}
			}
		}

		if child.Type == html.ElementNode && child.Data == "meta" {
			var content *string

			// Loop the meta element's attributes and detect what information
			// it carries. The loop is iterating backwards (from the last
			// attribute to the first), because the 'content' attribute is
			// likely placed after the 'name' attribute.
		recognizeMetaElement:
			for index, attr := len(child.Attr)-1, child.Attr[len(child.Attr)-1]; index >= 0; index, attr = index-1, child.Attr[index-1] {
				if attr.Namespace == "" && attr.Key == "content" {
					content = &attr.Val
					continue recognizeMetaElement
				}

				if attr.Namespace == "" && attr.Key == "charset" {
					webpage.Charset = attr.Val
					webpage.headNode.RemoveChild(child)
					continue processHead
				}

				if attr.Namespace == "" && attr.Key == "name" {
					if content == nil {
						for _, attr := range child.Attr {
							if attr.Namespace == "" && attr.Key == "content" {
								content = &attr.Val
								break
							}
						}

						if content == nil {
							continue recognizeMetaElement
						}
					}

					switch attr.Val {
					case "robots":
						webpage.Robots = *content

						webpage.headNode.RemoveChild(child)
						continue processHead
					case "viewport":
						webpage.Viewport = *content

						webpage.headNode.RemoveChild(child)
						continue processHead
					case "description":
						webpage.Description = *content

						webpage.headNode.RemoveChild(child)
						continue processHead
					case "author":
						webpage.Author = *content

						webpage.headNode.RemoveChild(child)
						continue processHead
					}
				}
			}
		}
	}

	return
}

func (webpage *Webpage) enrichHTML() {
	if webpage.noscriptNode.FirstChild != nil {
		webpage.injectNoscript()
	}

	if len(webpage.Meta) > 0 {
		webpage.injectMeta()
	}

	if len(webpage.Links) > 0 {
		webpage.injectLinks()
	}

	if len(webpage.Styles) > 0 {
		webpage.injectStyles()
	}

	if len(webpage.Scripts) > 0 {
		webpage.injectScripts()
	}

	if len(webpage.Title) > 0 {
		webpage.injectTitle()
	}

	return
}

func (webpage *Webpage) injectNoscript() {
	webpage.bodyNode.InsertBefore(webpage.noscriptNode, webpage.bodyNode.FirstChild)
}

func (webpage *Webpage) injectMeta() {
	allMeta := webpage.Meta

	addNamedContentMeta := func(name, content string) {
		if len(content) > 0 {
			allMeta = append(allMeta, []html.Attribute{
				{Key: "name", Val: name},
				{Key: "content", Val: content},
			})
		}
	}

	if len(webpage.Charset) > 0 {
		allMeta = append(allMeta, []html.Attribute{{
			Key: "charset",
			Val: webpage.Charset,
		}})
	}

	addNamedContentMeta("robots", webpage.Robots)
	addNamedContentMeta("viewport", webpage.Viewport)
	addNamedContentMeta("description", webpage.Description)
	addNamedContentMeta("author", webpage.Author)

	for _, meta := range allMeta {
		webpage.headNode.InsertBefore(&html.Node{
			Type: html.ElementNode,
			Data: "meta",
			Attr: meta,
		}, webpage.headNode.FirstChild)
	}
}

func (webpage *Webpage) injectLinks() {
	for _, link := range webpage.Links {
		webpage.headNode.InsertBefore(&html.Node{
			Type: html.ElementNode,
			Data: "link",
			Attr: link,
		}, webpage.headNode.FirstChild)
	}
}

func (webpage *Webpage) injectStyles() {
	for _, style := range webpage.Styles {
		styleNode := &html.Node{
			Type: html.ElementNode,
			Data: "style",
			Attr: style.Attributes,
		}

		styleNode.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: style.Content,
		})

		webpage.headNode.AppendChild(styleNode)
	}
}

func (webpage *Webpage) injectScripts() {
	for _, script := range webpage.Scripts {
		scriptNode := &html.Node{
			Type: html.ElementNode,
			Data: "script",
			Attr: script.Attributes,
		}

		if len(script.Src) > 0 {
			scriptNode.Attr = append(scriptNode.Attr, html.Attribute{
				Key: "src",
				Val: script.Src,
			})
		} else {
			scriptNode.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: script.Content,
			})
		}

		if script.Head {
			webpage.headNode.InsertBefore(scriptNode, webpage.headNode.FirstChild)
			continue
		}

		webpage.bodyNode.AppendChild(scriptNode)
	}
}

func (webpage *Webpage) injectTitle() {
	titleNode := webpage.findTitleNode()

	if len(webpage.Title) > 0 {
		titleNodeText := &html.Node{
			Type: html.TextNode,
			Data: webpage.Title,
		}

		if titleNode == nil {
			titleNode := &html.Node{
				Type: html.ElementNode,
				Data: "title",
				Attr: []html.Attribute{},
			}

			titleNode.AppendChild(titleNodeText)
			webpage.headNode.InsertBefore(titleNode, webpage.headNode.FirstChild)
		} else {
			for child := titleNode.FirstChild; child != nil; child = child.NextSibling {
				titleNode.RemoveChild(child)
			}

			titleNode.AppendChild(titleNodeText)
		}
	}
}

func (webpage *Webpage) findHeadNode() *html.Node {
	return findElement(webpage.rootNode, "head")
}

func (webpage *Webpage) findBodyNode() *html.Node {
	return findElement(webpage.rootNode, "body")
}

func (webpage *Webpage) findTitleNode() *html.Node {
	return findElement(webpage.headNode, "title")
}
