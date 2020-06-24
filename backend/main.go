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
	fmt.Println("server starting on port 5000")
	log.Fatal(fasthttp.ListenAndServe(":5000",
		router.Handler))
}
