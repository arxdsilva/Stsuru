package persisttest

import (
	"fmt"
	"log"
)

// FakeStorage is a interface that contains the fake functions that mongo implements
type FakeStorage interface {
	List() ([]string, error)
	Save(s string) error
	Remove(s string)
}

type FakeStore struct {
	Stored                      []string
	SaveErr, ListErr, RemoveErr error
}

func (f *FakeStore) Save(s string) error {
	if f.SaveErr == nil {
		f.Stored = append(f.Stored, s)
		return nil
	}
	f.SaveErr = fmt.Errorf("%s not saved", s)
	return f.SaveErr
}

func (f *FakeStore) List() ([]string, error) {
	return f.Stored, nil
}

func (f *FakeStore) Remove(s string) error {
	for i, e := range f.Stored {
		if s == e {
			f.Stored = append(f.Stored[:i], f.Stored[i+1:]...)
			return nil
		}
	}
	f.RemoveErr = fmt.Errorf("Could not remove %s", s)
	return f.RemoveErr
}

func checkError(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return err
}
