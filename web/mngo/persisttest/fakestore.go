package persisttest

import "log"

// FakeStorage ...
type FakeStorage struct {
	stored                                                  []string
	SaveErr, ListErr, RemoveErr, ExistHashErr, ExistLinkErr error
}

// Save ....
func (f *FakeStorage) Save(s string) error {
	if f.SaveErr == nil {
		f.stored = append(f.stored, s)
		return nil
	}
	return f.SaveErr
}

// List
func (f *FakeStorage) List() ([]string, error) {
	return f.stored, nil
}

func checkError(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return err
}
