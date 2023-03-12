package alox

import (
	"path"
	"strings"

	"golang.org/x/net/html"
)

func findNode(node *html.Node, predicate func(node *html.Node) bool) *html.Node {
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

func ShiftHead(value string) (head, tail string) {
	value = path.Clean("/" + value)
	splitIndex := strings.Index(value[1:], "/") + 1

	if splitIndex <= 0 {
		return value[1:], "/"
	}

	return value[1:splitIndex], value[splitIndex:]
}

func ShiftAndAssertHead(value string, assert func(head string) bool) (passed bool, tail string) {
	head, tail := ShiftHead(value)
	return assert(head), tail
}

func ShiftAndMatchHead(value string, head string) (matched bool, tail string) {
	return ShiftAndAssertHead(value, func(actualHead string) bool {
		return actualHead == head
	})
}

func AssertHead(value string, assert func(head string) bool) (passed bool) {
	head, _ := ShiftHead(value)
	return assert(head)
}

func MatchHead(value string, head string) (matched bool) {
	return AssertHead(value, func(actualHead string) bool {
		return actualHead == head
	})
}

func HasPrefix(value string, prefix string) bool {
	return strings.HasPrefix(path.Clean("/"+value), prefix)
}
