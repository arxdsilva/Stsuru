package main

import (
	"github.com/arxdsilva/Stsuru/web/persist"
	"github.com/arxdsilva/Stsuru/web/server"
)

func main() {
	// mongo := persist.FakeStore{
	// 	URL: "localhost",
	// }
	s := server.Server{Storage: &persist.FakeStore{}}
	s.Listen()
}

// usar no teste esta instancia para settar o FakeStorage como Storage do servidor, ao inves do mongo
// fake := persist.FakeStorage{}
