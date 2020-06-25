package main

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

// Este archivo tiene la función que recibe un dominio y hace uso de la
// librería o módulo d GO llamado GoQuery, que es similar a jQuery para
// poder hacer un poco de webscrapping y obtener el link (href) del logo

// Nota: Solo sirve con sitios web cuyo HTML está estructurado de tal manera
// que el logo se encuentre así:
// html -> head -> link y en el tag link tenga un atributo "rel" con valor
// "shortcut icon". Si lo encuentra, obtendrá el valor que está en el atributo
// "href" de dicho link, y allí debe estar el link del logo.

// Para otros sitios web que no lo tengan así, lamentablemente no lo arroja.
// Por ejemplo, Google no arroja logo, mientras que Netflix sí.

func GetHrefLinkLogo(url string) string {
	// Variable donde guardaremos el atributo de href
	var hrefDef string
	// Vamos a obtener el documento HTML haciendo goquery a la página
	doc, err := goquery.NewDocument("http://www." + url)
	// Si hubo un error obteniendolo lo notificamos
	if err != nil {
		log.Fatal(err)
	}
	// Vamos a ubicarnos en el arbol de nodos de HTML, donde hayamos
	// recorrido la secuencia html -> head -> link
	selection1 := doc.Find("html head link")
	// Ya posicionados en el lugar deseado, buscaremos entre todos los links
	// que hay el que tenga un atributo "rel" con valor de "shortcut icon",
	// y de existir obtendremos su atributo href.
	selection1.Each(func(_ int, selec *goquery.Selection) {
		rel, _ := selec.Attr("rel")
		if rel == "shortcut icon" {
			href, _ := selec.Attr("href")
			hrefDef = href
		}
	})
	// Vamos a retornar el valor de href
	return hrefDef
}
