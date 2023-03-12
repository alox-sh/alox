package webpage

import (
	"bytes"
	"testing"

	"golang.org/x/net/html"
)

func TestNewNoscriptNode(t *testing.T) {
	noscriptNode := newNoscriptNode()

	if noscriptNode.Type != html.ElementNode {
		t.Fatal("invalid (*html.Node).Type")
	}

	if noscriptNode.Data != "noscript" {
		t.Fatal("invalid (*html.Node).Data")
	}
}

func TestCloneNode(t *testing.T) {
	textA := "fooA"
	textB := "fooB"

	original := &html.Node{
		Type: html.ElementNode,
		Data: "div",
	}

	textNodeA := &html.Node{Type: html.TextNode, Data: textA}
	textNodeB := &html.Node{Type: html.TextNode, Data: textB}

	original.AppendChild(textNodeA)
	original.AppendChild(textNodeB)

	clone := cloneNode(original)

	if clone.Type != original.Type {
		t.Fatal("(*html.Node).Type does not match")
	}

	if clone.Data != original.Data {
		t.Fatal("(*html.Node).Data does not match")
	}

	isB := new(bool)
	*isB = false
	for child := clone.FirstChild; child != nil; child = child.NextSibling {
		if isB == nil {
			t.Fatal("expected exactly 2 childs of cloned node")
		}

		if child.Type != html.TextNode {
			t.Fatal("clone's child property (*html.Node).Type does not match")
		}

		if *isB {
			if child.Data != textB {
				t.Fatal("clone's child property (*html.Node).Data does not match")
			}

			child.Data = "barB"
			isB = nil
			continue
		}

		if child.Data != textA {
			t.Fatal("clone's child property (*html.Node).Data does not match")
		}

		child.Data = "barA"
		*isB = true
	}

	if textNodeA.Data != textA {
		t.Fatal("the clone is expected to be deep clone, but appears to be a shallow copy")
	}

	if textNodeB.Data != textB {
		t.Fatal("the clone is expected to be deep clone, but appears to be a shallow copy")
	}
}

func TestFindNode(t *testing.T) {
	document, err := html.Parse(bytes.NewReader([]byte(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>Foo bar</title>
			</head>
		</html>
	`)))
	if err != nil {
		t.Fatal(err)
	}

	titleNode := findNode(document, func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == "title"
	})

	if titleNode == nil {
		t.Fatal("expected to find 'title' element")
	}

	if titleNode.FirstChild == nil || titleNode.FirstChild.Type != html.TextNode || titleNode.FirstChild.Data != "Foo bar" {
		t.Fatal("invalid title element")
	}
}

func TestFindElement(t *testing.T) {
	document, err := html.Parse(bytes.NewReader([]byte(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>Foo bar</title>
			</head>
		</html>
	`)))
	if err != nil {
		t.Fatal(err)
	}

	titleNode := findElement(document, "title")

	if titleNode == nil {
		t.Fatal("expected to find 'title' element")
	}

	if titleNode.FirstChild == nil || titleNode.FirstChild.Type != html.TextNode || titleNode.FirstChild.Data != "Foo bar" {
		t.Fatal("invalid title element")
	}
}
