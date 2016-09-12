package mngo

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

var testCases = []struct {
	link  string
	hash  string
	isURL bool
}{
	{"abcdef", "", false},
	{"www.globo.com", "", false},
	{"www.notvalidurl.netscape", "", false},
	{"http://www.gorillatoolkit.org/pkg/mux", "70df8650c03c9fdfc959f04a64ecd956", true},
	{"https://mail.google.com/mail/u/0/#inbox", "2122c5656da3d86d77c08f7af48c0268", true},
}

func TestInsert(t *testing.T) {
	fmt.Print("Test Insert: ")
	r := httptest.NewRequest("POST", "/link/add", nil)
	w := httptest.NewRecorder()

	for _, test := range testCases {
		err := Insert(test.link, w, r)
		checkError(err)
		fmt.Print(".")
	}
	fmt.Println()
}

func TestGetAll(t *testing.T) {
	fmt.Print("Test GetAll: ")
	a, err := GetAll()
	checkError(err)
	if len(a) == 2 {
		fmt.Print("..")
	} else {
		fmt.Println()
		fmt.Printf("Array bigger than expected, len == %v", len(a))
	}
	fmt.Println()
}

func TestDelete(t *testing.T) {
	fmt.Print("Test Delete: ")
	for _, test := range testCases {
		err := Delete(test.hash)
		if err != nil {
			continue
		}
		fmt.Print(".")
	}
	fmt.Println()
}
