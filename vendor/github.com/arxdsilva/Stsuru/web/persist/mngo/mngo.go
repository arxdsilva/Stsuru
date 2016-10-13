package mngo

import (
	"github.com/arxdsilva/Stsuru/web/persist/data"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoStorage is the interface of CRUD methods
type MongoStorage struct {
	URL        string
	DB         string
	Collection string
}

// Save inputs a link into Mongo's DB
func (m *MongoStorage) Save(linkData *data.LinkData) error {
	s, err := mgo.Dial(m.URL)
	if err != nil {
		return err
	}
	defer s.Close()
	err = s.DB(m.DB).C(m.Collection).Insert(linkData)
	return err
}

// Remove removes a link from Mongo
func (m *MongoStorage) Remove(hash string) error {
	s, err := mgo.Dial(m.URL)
	if err != nil {
		return err
	}
	defer s.Close()

	c := s.DB(m.DB).C(m.Collection)
	err = c.Remove(bson.M{"hash": hash})
	return err
}

// FindHash finds an specific hash Stored on Mongo
func (m *MongoStorage) FindHash(hash string) (string, error) {
	s, err := mgo.Dial(m.URL)
	if err != nil {
		return "", err
	}
	defer s.Close()

	dbData := data.LinkData{}
	err = s.DB(m.DB).C(m.Collection).Find(bson.M{"hash": hash}).One(&dbData)
	return dbData.Link, err
}

// FindLink searches for an specific link inside Mongo
func (m *MongoStorage) FindLink(link string) (string, error) {
	s, err := mgo.Dial(m.URL)
	if err != nil {
		return "", err
	}
	defer s.Close()

	dbData := data.LinkData{}
	err = s.DB(m.DB).C(m.Collection).Find(bson.M{"link": link}).One(&dbData)
	return dbData.Link, err
}

// GetAll queries for all entries
func (m *MongoStorage) GetAll() ([]data.LinkData, error) {
	s, err := mgo.Dial(m.URL)
	if err != nil {
		return nil, err
	}
	defer s.Close()

	Data := []data.LinkData{}
	c := s.DB(m.DB).C(m.Collection)
	err = c.Find(bson.M{}).All(&Data)
	return Data, err
}

// CheckMultiple uses mongo to findout If a link was inserted twice
func (m *MongoStorage) CheckMultiple(link string, i int) (bool, error) {
	s, err := mgo.Dial(m.URL)
	if err != nil {
		return false, err
	}
	defer s.Close()

	dbNum := []data.LinkData{}
	err = s.DB(m.DB).C(m.Collection).Find(bson.M{"link": link}).All(&dbNum)
	if err != nil {
		return false, err
	}
	return len(dbNum) > i, nil
}
