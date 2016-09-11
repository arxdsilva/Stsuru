package mngo

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo is the interface of CRUD methods
type Mongo interface {
	Insert() error
	Delete() error
	FindOne() (string, error)
	FindAll() error
}

type lines struct {
	Link  string
	Short string
	Hash  string
}

// Insert ...
func Insert(link string) error {
	path := "http://localhost:8080/"
	// URL hashing
	linkShort, dbHash := hash(link, path)
	l := &lines{Link: link, Short: linkShort, Hash: dbHash}
	session, err := mgo.Dial("localhost")
	checkError(err)
	defer session.Close()
	err = session.DB("tsuru").C("links").Insert(l)
	return err
}

// Delete ...
func Delete(h string) error {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)
	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"hash": h})
	return err
}

// FindOne ...
func FindOne(s string) (string, error) {
	dbData := lines{}
	session, err := mgo.Dial("localhost")
	checkError(err)
	defer session.Close()
	err = session.DB("tsuru").C("links").Find(bson.M{"hash": s}).One(&dbData)
	return dbData.Link, err
}

// FindLink ...
func FindLink(s string) (string, error) {
	dbData := lines{}
	session, err := mgo.Dial("localhost")
	checkError(err)
	defer session.Close()
	err = session.DB("tsuru").C("links").Find(bson.M{"link": s}).One(&dbData)
	return dbData.Link, err
}

// GetAll ...
func GetAll() ([]lines, error) {
	Data := []lines{}
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)
	c := session.DB("tsuru").C("links")
	err = c.Find(bson.M{}).All(&Data)
	return Data, err
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}

func hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%s%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}
