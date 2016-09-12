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
	FindHash() (string, error)
	FindLink() (string, error)
	FindAll() error
}

type lines struct {
	Link  string
	Short string
	Hash  string
}

// Insert inputs a link into Mongo
func Insert(link string) error {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	path := "http://localhost:8080/"
	// URL hashing
	linkShort, dbHash := Hash(link, path)
	l := &lines{Link: link, Short: linkShort, Hash: dbHash}
	err = session.DB("tsuru").C("links").Insert(l)
	return err
}

// Delete removes a link from Mongo
func Delete(h string) error {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	c := session.DB("tsuru").C("links")
	err = c.Remove(bson.M{"hash": h})
	return err
}

// FindHash finds an specific hash Stored on Mongo
func FindHash(s string) (string, error) {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	dbData := lines{}
	err = session.DB("tsuru").C("links").Find(bson.M{"hash": s}).One(&dbData)
	return dbData.Link, err
}

// FindLink searches for an specific link inside Mongo
func FindLink(s string) (string, error) {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	dbData := lines{}
	err = session.DB("tsuru").C("links").Find(bson.M{"link": s}).One(&dbData)
	return dbData.Link, err
}

// GetAll queries for all entries
func GetAll() ([]lines, error) {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	Data := []lines{}
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

// Hash takes the URL, hashes & agregates the URL with your desired path
func Hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%sr/%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}
