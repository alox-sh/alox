package webpage

import (
	"golang.org/x/net/html"
)

func newNoscriptNode() *html.Node {
	return &html.Node{
		Type: html.ElementNode,
		Data: "noscript",
	}
}

func newHTMLDocument() *html.Node {
	return &html.Node{
		// Type: html.ElementNode,
		// Data: "noscript",
	}
}

func cloneNode(original *html.Node) (clone *html.Node) {
	attr := make([]html.Attribute, len(original.Attr))
	copy(attr, original.Attr)

	clone = &html.Node{
		Type:      original.Type,
		DataAtom:  original.DataAtom,
		Data:      original.Data,
		Namespace: original.Namespace,
		Attr:      attr,
	}

	for child := original.FirstChild; child != nil; child = child.NextSibling {
		clone.AppendChild(cloneNode(child))
	}

	return
}

func findNode(node *html.Node, predicate func(*html.Node) bool) *html.Node {
	if node == nil {
		return nil
	}

	if predicate(node) {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if foundNode := findNode(child, predicate); foundNode != nil {
			return foundNode
		}
	}

	return nil
}

func findElement(node *html.Node, name string) *html.Node {
	return findNode(node, func(node *html.Node) bool {
		if node.Type == html.ElementNode && node.Data == name {
			return true
		}

		return false
	})
}
