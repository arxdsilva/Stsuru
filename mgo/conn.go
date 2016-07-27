// This is a 'lib' to use mongoDB in other packages so we can manipulate the stored local DB
package mgo

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InitDB() *mgo.Session {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	return session
}

func InitSession(s *mgo.Session, db, collection string) *mgo.Collection {
	c := s.DB(db).C(collection)
	return c
}

func InsertP(c *mgo.Collection, n string, cp int, hp int, tp string) {
	err := c.Insert(
		// ??
		&Pokemon{n, cp, hp, tp},
	)
	if err != nil {
		panic(err)
	}
}

func ReadDB(c *mgo.Collection) {
	result := []Pokemon{}
	err := c.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func DeleteDB(c *mgo.Collection, n string) {
	_, err := c.RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}
}

func CloseDB(s *mgo.Session) {
	s.Close()
}
