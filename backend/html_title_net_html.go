package main

import (
	"io"

	"golang.org/x/net/html"
)

func esTitle(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "title"
}

func buscarRecursivamente(node *html.Node) (string, bool) {
	if esTitle(node) {
		return node.FirstChild.Data, true
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result, ok := buscarRecursivamente(child)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func ObtenerHTMLTitle(reader io.Reader) (string, bool) {
	htmlcontent, err := html.Parse(reader)
	if err != nil {
		panic("Fail to parse html")
	}

	return buscarRecursivamente(htmlcontent)
}
