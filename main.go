package main

import (
	"github.com/arxdsilva/Stsuru/web/persist"
	"github.com/arxdsilva/Stsuru/web/server"
)

func main() {
	mongo := persist.MongoStorage{
		URL:        "localhost",
		DB:         "tsuru",
		Collection: "links",
	}
	// usar no teste esta instancia para settar o FakeStorage como Storage do servidor, ao inves do mongo
	// fake := persisttest.FakeStorage{}
	server := server.Server{Storage: mongo}
	server.Listen()
}
