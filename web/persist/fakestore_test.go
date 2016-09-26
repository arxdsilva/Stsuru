package persist

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var expected = []Stored{
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
	for _, e := range expected {
		err := s.Save(e.Link)
		checkError(err)
		fmt.Print(".")
	}
	s.SaveErr = fmt.Errorf("not found")
	for _, e := range notexpected {
		err := s.Save(e.url)
		if err == fmt.Errorf("%s not saved", e.url) {
			fmt.Print(".")
		}
	}
	fmt.Println()
}

func TestList(t *testing.T) {
	fmt.Print("Testing List: ")
	s := FakeStore{
		Stored: expected,
	}
	list, err := s.List()
	checkError(err)
	if !reflect.DeepEqual(expected, list) {
		log.Panicf("List %v is not equal to list %v", list, expected)
	}
	fmt.Println(".")
}

func TestExists(t *testing.T) {
	fmt.Print("Testing Exists: ")
	s := FakeStore{
		Stored: expected,
	}
	for _, e := range expected {
		result := s.Exists(e.Link)
		if result == true {
			fmt.Print(".")
			continue
		}
		log.Panicf("Element %s could not be found on slice %v", e.Link, expected)
	}
	for _, e := range notexpected {
		r := s.Exists(e.url)
		if r == e.value {
			fmt.Print(".")
			continue
		}
		log.Panicf("Element %s should not be found on slice %v", e.url, expected)
	}
	fmt.Println()
}

func TestFindHash(t *testing.T) {
	fmt.Print("Testing FindHash: ")
	s := FakeStore{
		Stored: expected,
	}
	for _, e := range expected {
		_, err := s.FindHash(e.Hash)
		if err != nil {
			log.Panicf("Element %s was not found in %v", e.Hash, expected)
		}
		fmt.Print(".")
	}
	for _, e := range notexpected {
		_, err := s.FindHash(e.url)
		if err != nil {
			fmt.Print(".")
		}
	}
	fmt.Println()
}

func TestRemove(t *testing.T) {
	fmt.Print("Testing Remove: ")
	s := FakeStore{
		Stored: expected,
	}
	for _, e := range expected {
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