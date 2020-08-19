package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// El programa principal del servidor.

func main() {
	// Se instancia el API Router de fasthttprouter
	router := fasthttprouter.New()
	// Se declara el endpoint 2, con método GET, donde se obtienen los dominios consultados (obtenidos de la BD)
	router.GET("/queried/domains", GetQueriedDomains)
	// Se declara el endpoint 1, con método POST, que recibe un dominio en la URL (domain) y
	// realiza todas las solicitudes y procesos para obtener la info de los servers e interactuar
	// con la BD
	router.POST("/info/servers/:domain", PostDomainAndGetInfo)
	// Se le ordena al server fasthttp que empieze a escuchar en el puerto 5000 y use
	// los handlers que están attacheados al router
	log.Fatal(fasthttp.ListenAndServe(":5000",
		router.Handler))
	// Se avisa que el server empieza a escuchar en el puerto 5000
	fmt.Print("server starting on port 5000")
	fmt.Print("Backend is ready on port 5000 of localhost")
}
