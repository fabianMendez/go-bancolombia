package main

import "golang.org/x/net/html"

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

func getElementByID(node *html.Node, id string) *html.Node {
	return getElementMatching(node, func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			elmID := getAttribute(n, "id")
			return elmID == id
		}

		return false
	})
}
