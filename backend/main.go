package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	router.GET("/queried/domains", GetQueriedDomains)
	router.POST("/info/servers/:domain", PostDomainAndGetInfo)
	fmt.Println("server starting on localhost:5000")
	log.Fatal(fasthttp.ListenAndServe("localhost:5000",
		router.Handler))
}
