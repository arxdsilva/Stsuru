package persist

import (
	"fmt"
	"log"
)

// Stored is the fake DB
type Stored struct {
	Link, LinkShort, Hash string
}

// FakeStore is
type FakeStore struct {
	URL                         string
	Stored                      []Stored
	SaveErr, ListErr, RemoveErr error
}

// Save ...
func (f *FakeStore) Save(link, linkShort, dbHash string) error {
	n := Stored{
		Link: link,
		Hash: dbHash,
	}
	if f.SaveErr == nil {
		f.Stored = append(f.Stored, n)
		return nil
	}
	f.SaveErr = fmt.Errorf("%s not saved", link)
	return f.SaveErr
}

// List ...
func (f *FakeStore) List() ([]Stored, error) {
	return f.Stored, nil
}

// Remove ...
func (f *FakeStore) Remove(dbHash string) error {
	for i, e := range f.Stored {
		if dbHash == e.Hash {
			f.Stored = append(f.Stored[:i], f.Stored[i+1:]...)
			return nil
		}
	}
	f.RemoveErr = fmt.Errorf("Could not remove %s", dbHash)
	return f.RemoveErr
}

// Exists ...
func (f *FakeStore) Exists(s string) bool {
	for _, e := range f.Stored {
		if s == e.Link {
			return true
		}
	}
	return false
}

// FindHash ...
func (f *FakeStore) FindHash(s string) (string, error) {
	for _, e := range f.Stored {
		if e.Hash == s {
			return e.Link, fmt.Errorf("Found")
		}
	}
	return s, nil
}

func checkError(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
