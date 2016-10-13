package persist

import (
	"fmt"

	"github.com/arxdsilva/Stsuru/web/persist/data"
)

// Stored is the fake DB
type Stored struct {
	Link, LinkShort, Hash string
}

// FakeStore stores the given data into FakeStore.Stored so user can use It as
// a fake DB
type FakeStore struct {
	URL                         string
	Stored                      []Stored
	SaveErr, ListErr, RemoveErr error
}

// Save ...
func (f *FakeStore) Save(linkData *data.LinkData) error {
	n := Stored{
		Link:      linkData.Link,
		LinkShort: linkData.Short,
		Hash:      linkData.Hash,
	}
	if f.SaveErr != nil {
		f.SaveErr = fmt.Errorf("%s not saved", n.Link)
		return f.SaveErr
	}
	f.Stored = append(f.Stored, n)
	return nil
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
			return e.Link, nil
		}
	}
	return s, fmt.Errorf("Not Found")
}
