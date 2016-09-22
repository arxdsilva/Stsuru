package persist

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var Expected = []Stored{
	{"http://localhost:8080/", "9825c2a542dd888e55b9b0e06b04f672"},
	{"http://science.nasa.gov/", "af13587359208048616bfedcb3b4dbdc"},
}
var notexpected = []struct {
	url   string
	value bool
}{
	{"ssssssscience.nasa.gov/", false},
	{"https://mail.google.com/mail/u/1/#inbox", false},
}

func TestSave(t *testing.T) {
	fmt.Print("Testing Save: ")
	s := FakeStore{}
	for _, e := range Expected {
		err := s.Save(e.Link)
		checkError(err)
		fmt.Print(".")
	}
	fmt.Println()
}

func TestList(t *testing.T) {
	fmt.Print("Testing List: ")
	s := FakeStore{
		Stored: Expected,
	}
	list, err := s.List()
	checkError(err)
	if !reflect.DeepEqual(Expected, list) {
		log.Panicf("List %v is not equal to list %v", list, Expected)
	}
	fmt.Println(".")
}

func TestExists(t *testing.T) {
	fmt.Print("Testing Exists: ")
	s := FakeStore{
		Stored: Expected,
	}
	for _, e := range Expected {
		result := s.Exists(e.Link)
		if result == true {
			fmt.Print(".")
			continue
		}
		log.Panicf("Element %s could not be found on slice %v", e.Link, Expected)
	}
	for _, e := range notexpected {
		r := s.Exists(e.url)
		if r == e.value {
			fmt.Print(".")
			continue
		}
		log.Panicf("Element %s should not be found on slice %v", e.url, Expected)
	}
	fmt.Println()
}

func TestFindHash(t *testing.T) {
	fmt.Print("Testing FindHash: ")
	s := FakeStore{
		Stored: Expected,
	}
	for _, e := range Expected {
		_, err := s.FindHash(e.Hash)
		if err != nil {
			log.Panicf("Element %s was not found in %v", e.Hash, Expected)
		}
		fmt.Print(".")
	}
	fmt.Println()
}

func TestRemove(t *testing.T) {
	fmt.Print("Testing Remove: ")
	s := FakeStore{
		Stored: Expected,
	}
	for _, e := range Expected {
		err := s.Remove(e.Hash)
		checkError(err)
		fmt.Print(".")
	}
	for _, e := range notexpected {
		err := s.Remove(e.url)
		if err != nil {
			fmt.Print(".")
			continue
		}
		log.Panicf("Expected %s and received %v", "not found", err)
	}
	fmt.Println()
}
