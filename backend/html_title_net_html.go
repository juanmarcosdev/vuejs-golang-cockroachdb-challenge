package main

import (
	"io"

	"golang.org/x/net/html"
)

// Este archivo utiliza las funciones aquí escritas apra poder obtener
// el title del HTML (lo busca de manera recursiva, usando el paquete
// modulo de GO de HTML).

// Esta función verifica si de un nodo HTML que recibe corresponde
// al title
func esTitle(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "title"
}

// Esta función recibe un nodo HTML y busca de manera recursiva
// hasta encontrar el nodo que corresponda al title
func buscarRecursivamente(node *html.Node) (string, bool) {
	// Por cada nodo HTML que recibe verifica si es title, si es (caso base),
	// sale de la función y lo retorna
	if esTitle(node) {
		return node.FirstChild.Data, true
	}
	// De no serlo (caso recursivo), se encarga de buscar
	// recursivamente por cada nodo hijo en el HTML
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result, ok := buscarRecursivamente(child)
		if ok {
			return result, ok
		}
	}
	// Si al final no hay nodo HTML retorna string vacío
	return "", false
}

// Esta función ejecuta la de arriba hasta que lo encuentra, y retorna
// dicho title
func ObtenerHTMLTitle(reader io.Reader) (string, bool) {
	// Hace el parse-HTML al body de la respuesta que recibe
	htmlcontent, err := html.Parse(reader)
	// Si hay errores en el parseo se notifican
	if err != nil {
		panic("No se pudo parsear HTML")
	}
	// Realiza búsqueda recursiva
	return buscarRecursivamente(htmlcontent)
}
