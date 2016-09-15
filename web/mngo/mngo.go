package mngo

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"

	"github.com/asaskevich/govalidator"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage is the interface of CRUD methods
type Storage interface {
	Insert(string) error
	Delete(string) error
	FindHash(string) error
	FindLink(string) error
	FindAll() error
}

// Mongo ...
type Mongo struct{}

type LinkData struct {
	Link  string
	Short string
	Hash  string
}

// Insert inputs a link into Mongo
func Insert(link string) error {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)
	v := validateURL(link)
	if v == true {
		_, err = FindLink(link)
		if err == nil {
			return err
		}
		path := "http://localhost:8080/"
		linkShort, dbHash := Hash(link, path)
		l := &LinkData{Link: link, Short: linkShort, Hash: dbHash}
		err = session.DB("tsuru").C("links").Insert(l)
		return err
	}
	return nil
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

	dbData := LinkData{}
	err = session.DB("tsuru").C("links").Find(bson.M{"hash": s}).One(&dbData)
	return dbData.Link, err
}

// FindLink searches for an specific link inside Mongo
func FindLink(s string) (string, error) {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	dbData := LinkData{}
	err = session.DB("tsuru").C("links").Find(bson.M{"link": s}).One(&dbData)
	return dbData.Link, err
}

// GetAll queries for all entries
func GetAll() ([]LinkData, error) {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	Data := []LinkData{}
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

// Hash creates & returns a link with the hashed URL and the URL hash
func Hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%s%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}

func validateURL(l string) bool {
	isURL := govalidator.IsURL(l)
	validURL := govalidator.IsRequestURL(l)
	if isURL == false || validURL == false {
		return false
	}
	return true
}

// CheckMultiple ...
func CheckMultiple(s string, i int) bool {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	dbNum := []LinkData{}
	err = session.DB("tsuru").C("links").Find(bson.M{"link": s}).All(&dbNum)
	checkError(err)
	if len(dbNum) > i {
		return true
	}
	return false
}
