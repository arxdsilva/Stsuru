package persisttest

import (
	"fmt"
	"log"
)

// FakeStorage commentary
type FakeStorage struct {
	Stored                      []string
	SaveErr, ListErr, RemoveErr error
}

// Save ....
func (f *FakeStorage) Save(s string) error {
	if f.SaveErr == nil {
		f.Stored = append(f.Stored, s)
		return nil
	}
	f.SaveErr = fmt.Errorf("%s not saved", s)
	return f.SaveErr
}

// List ...
func (f *FakeStorage) List() ([]string, error) {
	return f.Stored, nil
}

// Remove ...
func (f *FakeStorage) Remove(s string) error {
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
