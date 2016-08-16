// This is a 'lib' to use mongoDB in other packages so we can manipulate the stored local DB
package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	session := InitDB()
	c := InitSession(session, "nintendo", "pokemons")
	InsertP(c, "charmander", 10, 100, "fire")
	ReadDB(c)
	DeleteDB(c, "charmander")
	ReadDB(c)
	CloseDB(session)
}

type Pokemon struct {
	Name string
	CP   int
	HP   int
	Type string
}

// InitDB ,,,,,
func InitDB() *mgo.Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	return session
}

// InitSession ,,,,,
func InitSession(s *mgo.Session, db, collection string) *mgo.Collection {
	c := s.DB(db).C(collection)
	return c
}

// InsertP ,,,,,
func InsertP(c *mgo.Collection, n string, cp int, hp int, tp string) {
	err := c.Insert(
		&Pokemon{n, cp, hp, tp},
	)
	if err != nil {
		panic(err)
	}
}

// ReadDB ,,,,,
func ReadDB(c *mgo.Collection) {
	result := []Pokemon{}
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// DeleteDB ,,,,,
func DeleteDB(c *mgo.Collection, n string) {
	_, err := c.RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}
}

// CloseDB ,,,,,
func CloseDB(s *mgo.Session) {
	s.Close()
}
