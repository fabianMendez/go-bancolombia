package main

import (
	"net/url"

	"golang.org/x/net/html"
)

func getAttribute(n *html.Node, attrKey string) string {
	for _, a := range n.Attr {
		if a.Key == attrKey {
			return a.Val
		}
	}

	return ""
}

func getElementMatching(node *html.Node, fn func(*html.Node) bool) *html.Node {
	if fn(node) {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if elm := getElementMatching(child, fn); elm != nil {
			return elm
		}
	}

	return nil
}

func getAllElemenstMatching(node *html.Node, fn func(*html.Node) bool) []*html.Node {
	var elms []*html.Node

	if fn(node) {
		elms = append(elms, node)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		elms = append(elms, getAllElemenstMatching(child, fn)...)
	}

	return elms
}

func getElementByID(node *html.Node, id string) *html.Node {
	return getElementMatching(node, func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			elmID := getAttribute(n, "id")
			return elmID == id
		}

		return false
	})
}

func hasTag(node *html.Node, tag string) bool {
	return node.Type == html.ElementNode && node.Data == tag
}

func getAllElemenstByTag(node *html.Node, tag string) []*html.Node {
	return getAllElemenstMatching(node, func(n *html.Node) bool {
		return hasTag(n, tag)
	})
}

func getInnerText(n *html.Node) string {
	str := ""

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			str += c.Data
		} else if c.Type == html.ElementNode && c.Data == "br" {
			str += "\n"
			// } else {
		} else if c.Type == html.ElementNode && c.Data != "script" {
			str += getInnerText(c)
		}
	}

	return str
}
func parseFormFields(form *html.Node) url.Values {
	fields := url.Values{}

	inputs := getAllElemenstByTag(form, "input")

	for _, input := range inputs {
		name := getAttribute(input, "name")
		value := getAttribute(input, "value")
		fields.Add(name, value)
	}

	return fields
}
