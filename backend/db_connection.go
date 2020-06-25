package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DBConnector struct {
	connection *sql.DB
}

// Este archivo corresponde a la función que crea un nuevo conector de base de datos SQL,
// según se le pase un usuario, un nombre de la base de datos y un puerto. Automáticamente se
// conecta al host que tenga por nombre <database> (Docker será capaz de resolver el nombre, ya que
// estarán conectados en la misma red gracias a Docker Compose).

// Realmente esta función es muy sencilla: hace uso de la librería sql embebida en Go y el comando
// Open para abrir una nueva conexión a la base de datos que usará el servidor para escribir los dominios
// y su información en el primer endpoint, y obtener los dominios consultados en el segundo endpoint.

func NewDBConnector(user string, dbname string, port string) (*DBConnector, error) {
	newConnec, err := sql.Open("postgres", "host = database user = "+user+" dbname = "+dbname+" sslmode=disable port = "+port)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &DBConnector{
		connection: newConnec,
	}, nil
}
