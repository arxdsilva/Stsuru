package main

import (
	"github.com/arxdsilva/Stsuru/web/persist"
	"github.com/arxdsilva/Stsuru/web/server"
)

func main() {
	s := server.Server{Storage: &persist.FakeStore{}}
	s.Listen()
}
