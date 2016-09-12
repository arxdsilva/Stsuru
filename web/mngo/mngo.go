package mngo

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"

	"github.com/asaskevich/govalidator"

	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo is the interface of CRUD methods
type Mongo interface {
	Insert()
	Delete()
	FindHash()
	FindLink()
	FindAll()
}

type lines struct {
	Link  string
	Short string
	Hash  string
}

// Insert inputs a link into Mongo
func Insert(link string, w http.ResponseWriter, r *http.Request) error {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	valid := validateURL(link)
	if valid == true {
		_, err = FindLink(link)
		if err == nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}
	path := "http://localhost:8080/"
	linkShort, dbHash := hash(link, path)
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

func hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%sr/%x", path, hash)
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

func checkMultiple(s string) ([]lines, bool) {
	session, err := mgo.Dial("localhost")
	defer session.Close()
	checkError(err)

	dbNum := []lines{}
	err = session.DB("tsuru").C("links").Find(bson.M{"link": s}).All(&dbNum)
	checkError(err)
	if len(dbNum) > 1 {
		return dbNum, true
	}
	return dbNum, false
}
