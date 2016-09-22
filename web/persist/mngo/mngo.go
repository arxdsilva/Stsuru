package mngo

import (
	"log"

	"github.com/arxdsilva/Stsuru/web/persist"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoStorage is the interface of CRUD methods
type MongoStorage struct {
	URL        string
	DB         string
	Collection string
}

// LinkData holds the structure that is used by mongo to insert data to DB
type LinkData struct {
	Link  string
	Short string
	Hash  string
}

// Save inputs a link into Mongo's DB
func (m *MongoStorage) Save(link string) error {
	s, err := mgo.Dial(m.URL)
	defer s.Close()
	checkError(err)
	v := persist.ValidateURL(link)
	if v == true {
		_, err = m.FindLink(link)
		if err == nil {
			return err
		}
		linkShort, dbHash := persist.Hash(link, m.URL)
		l := &LinkData{Link: link, Short: linkShort, Hash: dbHash}
		err = s.DB(m.DB).C(m.Collection).Insert(l)
		return err
	}
	return nil
}

// Remove removes a link from Mongo
func (m *MongoStorage) Remove(h string) error {
	s, err := mgo.Dial(m.URL)
	defer s.Close()
	checkError(err)

	c := s.DB(m.DB).C(m.Collection)
	err = c.Remove(bson.M{"hash": h})
	return err
}

// FindHash finds an specific hash Stored on Mongo
func (m *MongoStorage) FindHash(hash string) (string, error) {
	s, err := mgo.Dial(m.URL)
	defer s.Close()
	checkError(err)

	dbData := LinkData{}
	err = s.DB(m.DB).C(m.Collection).Find(bson.M{"hash": hash}).One(&dbData)
	return dbData.Link, err
}

// FindLink searches for an specific link inside Mongo
func (m *MongoStorage) FindLink(link string) (string, error) {
	s, err := mgo.Dial(m.URL)
	defer s.Close()
	checkError(err)

	dbData := LinkData{}
	err = s.DB(m.DB).C(m.Collection).Find(bson.M{"link": link}).One(&dbData)
	return dbData.Link, err
}

// GetAll queries for all entries
func (m *MongoStorage) GetAll() ([]LinkData, error) {
	s, err := mgo.Dial(m.URL)
	defer s.Close()
	checkError(err)

	Data := []LinkData{}
	c := s.DB(m.DB).C(m.Collection)
	err = c.Find(bson.M{}).All(&Data)
	return Data, err
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
	return
}

// CheckMultiple uses mongo to findout If a link was inserted twice
func (m *MongoStorage) CheckMultiple(link string, i int) bool {
	s, err := mgo.Dial(m.URL)
	defer s.Close()
	checkError(err)

	dbNum := []LinkData{}
	err = s.DB(m.DB).C(m.Collection).Find(bson.M{"link": link}).All(&dbNum)
	checkError(err)
	if len(dbNum) > i {
		return true
	}
	return false
}
