package mgo

import (
	"gopkg.in/mgo.v2"
)

// NewConn returns a new connection with database
func NewConn() {
	conn, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	return conn
}

func ()  {

}
