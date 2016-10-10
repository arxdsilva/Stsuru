package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/arxdsilva/Stsuru/web/persist"
	"github.com/gorilla/mux"
)

var testCases = []struct {
	name     string
	expected int
}{
	{"https://www.youtube.com/watch?v=WC5FdFlUcl0&list=PL2LC4RhdHOKe3y0pDyNBMF9ztULgUzKOq&index=1", http.StatusFound},
	{"http://science.nasa.gov/", http.StatusFound},
	{"https://godoc.org/gopkg.in/mgo.v2", http.StatusFound},
	{"http://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go", http.StatusNotFound},
}
var s = Server{Storage: &persist.FakeStore{}}

func TestHome(t *testing.T) {
	dir := "../../"
	err := os.Chdir(dir)
	if err != nil {
		t.Fatalf("unexpected error changing dir: %v", err)
	}
	fmt.Print("Testing Home: ")
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("unexpected error new request: %v", err)
	}
	w := httptest.NewRecorder()
	s.Home(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
		return
	}
	fmt.Println(".")
}

func TestAddLink(t *testing.T) {
	var testURLs = []struct {
		name   string
		expect int
	}{
		{"notalink", http.StatusNotModified},
		{"notavalidurl.com", http.StatusNotModified},
		{"https://www.youtube.com/watch?v=WC5FdFlUcl0&list=PL2LC4RhdHOKe3y0pDyNBMF9ztULgUzKOq&index=1", http.StatusFound},
		{"http://science.nasa.gov/", http.StatusFound},
		{"multiple.dots.not.valid.url", http.StatusNotModified},
		{"https://godoc.org/gopkg.in/mgo.v2", http.StatusFound},
	}
	v := url.Values{}

	fmt.Print("Test Add Link: ")
	for _, test := range testURLs {
		v.Set("user_link", test.name)
		tf := strings.NewReader(v.Encode())
		r, err := http.NewRequest("POST", "/link/add", tf)
		if err != nil {
			t.Fatalf("unexpected error new request: %v", err)
		}
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		s.AddLink(w, r)

		if w.Code != test.expect {
			t.Errorf("Request %s on Home page returned %v instead of %v", test.name, w.Code, test.expect)
		}
		fmt.Print(".")
	}
	i := 3
	stored, err := s.Storage.List()
	if err != nil {
		t.Fatalf("unexpected error storage list: %v", err)
	}
	if len(stored) != i {
		t.Errorf("Storage has inserted multiple equal URLs")
	}
	fmt.Print(". ")
	fmt.Println()
}

func TestRedirect(t *testing.T) {
	fmt.Print("Test Redirect: ")
	var link string
	stored, _ := s.Storage.List()
	for i, test := range testCases {
		for _, obj := range stored {
			if test.name == obj.Link {
				link = obj.Hash
			}
			continue
		}
		if i == 3 {
			link = "llll"
		}
		path := "/r/"
		n := addPath(path, link)
		r, err := http.NewRequest("GET", n, nil)
		if err != nil {
			t.Fatalf("unexpected error new request: %v", err)
		}
		r.Header.Set("Content-Type", "text/html")
		r.Header.Add("Accept", "text/html")
		r.Header.Set("Accept", "application/xhtml+xml")
		w := httptest.NewRecorder()

		m := mux.NewRouter()
		m.HandleFunc("/r/{id}", s.Redirect)
		m.ServeHTTP(w, r)

		if w.Code != test.expected {
			fmt.Println(n)
			fmt.Println(test.name)
			t.Errorf("\nLink %s got %v instead of %v\n", link, w.Code, test.expected)
			continue
		}
		fmt.Print("* ")
	}
	fmt.Println()
}

func TestRemoveLink(t *testing.T) {
	fmt.Print("Test Removing Links: ")
	path := "/l/r/"
	stored, _ := s.Storage.List()
	for _, test := range stored {
		pathed := addPath(path, test.Hash)
		r, err := http.NewRequest("GET", pathed, nil)
		if err != nil {
			t.Fatalf("unexpected error new request: %v", err)
		}
		r.Header.Set("Content-Type", "text/html")
		r.Header.Add("Accept", "text/html")
		r.Header.Set("Accept", "application/xhtml+xml")
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc("/l/r/{id}", s.RemoveLink)
		m.ServeHTTP(w, r)
		_, err = s.Storage.FindHash(test.Hash)
		if err == nil {
			fmt.Printf("\n%s not expected on Storage", test.Hash)
			continue
		}
		fmt.Print("x ")
	}
	fmt.Println()
}

func addPath(path, hash string) string {
	return fmt.Sprintf("%s%s", path, hash)
}
