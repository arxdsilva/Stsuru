package persist

import (
	"fmt"
	"log"
)

type Stored struct {
	Link, Hash string
}

type FakeStore struct {
	URL                         string
	DB                          string
	Collection                  string
	Stored                      []Stored
	SaveErr, ListErr, RemoveErr error
}

// Storage is a interface that contains the fake functions that mongo implements
type Storage interface {
	List() ([]Stored, error)
	Save(s string) error
	Remove(s string) error
	Exists(s string) error
	FindHash(s string) (string, error)
}

func (f *FakeStore) Save(s string) error {
	l, dbD := Hash(s, f.URL)
	n := Stored{
		Link: l,
		Hash: dbD,
	}
	if f.SaveErr == nil {
		f.Stored = append(f.Stored, n)
		return nil
	}
	f.SaveErr = fmt.Errorf("%s not saved", s)
	return f.SaveErr
}

func (f *FakeStore) List() ([]Stored, error) {
	return f.Stored, nil
}

func (f *FakeStore) Remove(s string) error {
	for i, e := range f.Stored {
		if s == e.Hash {
			f.Stored = append(f.Stored[:i], f.Stored[i+1:]...)
			return nil
		}
	}
	f.RemoveErr = fmt.Errorf("Could not remove %s", s)
	return f.RemoveErr
}

func (f *FakeStore) Exists(s string) bool {
	for _, e := range f.Stored {
		if s == e.Link {
			return true
		}
	}
	return false
}

func (f *FakeStore) FindHash(s string) (string, error) {
	for _, e := range f.Stored {
		if e.Hash == s {
			return e.Link, nil
		}
	}
	return s, fmt.Errorf("%s not found", s)
}

func checkError(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
