package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DBConnector struct {
	connection *sql.DB
}

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
