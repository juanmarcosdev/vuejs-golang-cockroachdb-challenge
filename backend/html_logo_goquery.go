package main

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

func GetHrefLinkLogo(url string) string {
	var hrefDef string
	doc, err := goquery.NewDocument("http://www." + url)
	if err != nil {
		log.Fatal(err)
	}
	selection1 := doc.Find("html head link")
	selection1.Each(func(_ int, selec *goquery.Selection) {
		rel, _ := selec.Attr("rel")
		if rel == "shortcut icon" {
			href, _ := selec.Attr("href")
			hrefDef = href
		}
	})
	return hrefDef
}
