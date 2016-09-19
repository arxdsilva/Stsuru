package persisttest

import (
	"fmt"
	"reflect"
	"testing"
)

var expected = []string{"https://www.globo.com", "https://www.google.com"}

func TestSave(t *testing.T) {
	fmt.Print("Testing Save: ")
	s := FakeStorage{}
	for _, e := range expected {
		err := s.Save(e)
		checkError(err)
		fmt.Print(".")
	}
	fmt.Println()
}

func TestList(t *testing.T) {
	fmt.Print("Testing List: ")
	s := FakeStorage{}
	list, err := s.List()
	checkError(err)
	if !reflect.DeepEqual(expected, list) {
		fmt.Printf("List %v is not equal to list %v", list, expected)
	}
	fmt.Println(".")
}
