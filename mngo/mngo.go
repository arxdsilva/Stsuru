package mngo

import (
  "gopkg.in/mgo.v2/bson"
  "github.com/arxdsilva/Stsuru/handlers"
)

// Insert ...
func Insert(l *Lines) error {
  session, err := mgo.Dial("localhost")
  handlers.CheckError(err)
  defer session.Close()
	err = session.DB("tsuru").C("links").Insert(l)
	return err
}

// Delete ...
func Delete(h string) error {
  session, err := mgo.Dial("localhost")
	defer session.Close()
	handlers.CheckError(err)
	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"hash": h})
	return err
}

// FindOne ...
func FindOne(s string) (string, error) {
  dbData := handlers.Lines{}
	session, err := mgo.Dial("localhost")
	checkError(err)
	defer session.Close()
	err = session.DB("tsuru").C("links").Find(bson.M{"hash": dbHash}).One(&dbData)
	return dbData.Link, err
}

// FindAll ...
func FindAll(d *Data) {
  session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)
	c := session.DB("tsuru").C("links")
	err = c.Find(bson.M{}).All(&Data)
	return err
}
