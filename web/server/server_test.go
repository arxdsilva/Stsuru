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
	"github.com/arxdsilva/Stsuru/web/persist/mngo"
	"github.com/gorilla/mux"
)

var testCases = []struct {
	name     string
	hash     string
	expected int
}{
	{"http://localhost:8080/", "9825c2a542dd888e55b9b0e06b04f672", http.StatusFound},
	{"http://science.nasa.gov/", "af13587359208048616bfedcb3b4dbdc", http.StatusFound},
	{"https://godoc.org/gopkg.in/mgo.v2", "b5cfe5dac82a4a8af7a505891cd91729", http.StatusFound},
	{"http://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go", "", http.StatusNotFound},
}
var s = Server{Storage: &persist.FakeStore{}}

func TestHome(t *testing.T) {
	dir := "../../"
	err := os.Chdir(dir)
	checkError(err)
	fmt.Print("Testing Home: ")
	r := httptest.NewRequest("GET", "/", nil)
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
		{"", http.StatusNotModified},
		{"notalink", http.StatusNotModified},
		{"notavalidurl.com", http.StatusNotModified},
		{"http://localhost:8080/", http.StatusFound},
		{"http://science.nasa.gov/", http.StatusFound},
		{"multiple.dots.not.valid.url", http.StatusNotModified},
		{"https://godoc.org/gopkg.in/mgo.v2", http.StatusFound},
		{"https://godoc.org/gopkg.in/mgo.v2", http.StatusNotModified},
	}
	v := url.Values{}

	fmt.Print("Test Add Link: ")
	for _, test := range testURLs {
		v.Set("user_link", test.name)
		tf := strings.NewReader(v.Encode())
		r := httptest.NewRequest("POST", "/link/add", tf)
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
	checkError(err)
	if len(stored) != i {
		t.Errorf("Storage has inserted multiple equal URLs")
	}
	fmt.Print(". ")
	fmt.Println()
}

func TestRedirect(t *testing.T) {
	fmt.Print("Test Redirect: ")
	for _, test := range testCases {
		link := test.name
		path := "/r/"
		n, _ := hash(link, path)
		r := httptest.NewRequest("GET", n, nil)
		r.Header.Set("Content-Type", "text/html")
		r.Header.Add("Accept", "text/html")
		r.Header.Set("Accept", "application/xhtml+xml")
		w := httptest.NewRecorder()

		m := mux.NewRouter()
		m.HandleFunc("/r/{id}", s.Redirect)
		m.ServeHTTP(w, r)
		if w.Code != test.expected {
			t.Errorf("\nLink %s got %v instead of %v\n", link, w.Code, test.expected)
			continue
		}
		fmt.Print("* ")
	}
	fmt.Println()
}

func TestRemoveLink(t *testing.T) {
	fmt.Print("Test Removing Links: ")
	for _, test := range testCases {
		path := "http://tsu.ru:8080/l/r/"
		link := test.name
		n, dbHash := hash(link, path)
		ngo := mngo.MongoStorage{}

		r := httptest.NewRequest("GET", n, nil)
		r.Header.Set("Content-Type", "text/html")
		r.Header.Add("Accept", "text/html")
		r.Header.Set("Accept", "application/xhtml+xml")
		w := httptest.NewRecorder()
		m := mux.NewRouter()
		m.HandleFunc("/l/r/{id}", s.RemoveLink)
		m.ServeHTTP(w, r)

		_, err := ngo.FindHash(dbHash)
		if err == nil {
			fmt.Printf("\n%s not expected on Mongo", dbHash)
			continue
		}
		fmt.Print("x ")
	}
	fmt.Println()
}
