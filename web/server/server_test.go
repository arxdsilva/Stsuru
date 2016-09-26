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

var mongo = persist.MongoStorage{
	URL:        "localhost",
	DB:         "tsuru",
	Collection: "links",
}

var svr = Server{Storage: mongo}

var testCases = []struct {
	name string
	hash string
}{
	{"http://localhost:8080/", "9825c2a542dd888e55b9b0e06b04f672"},
	{"http://science.nasa.gov/", "af13587359208048616bfedcb3b4dbdc"},
	{"https://godoc.org/gopkg.in/mgo.v2", "b5cfe5dac82a4a8af7a505891cd91729"},
}
var s = Server{}

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
		expect bool
	}{
		{"", false},
		{"notalink", false},
		{"notavalidurl.com", false},
		{"http://localhost:8080/", true},
		{"http://science.nasa.gov/", true},
		{"multiple.dots.not.valid.url", false},
		{"https://godoc.org/gopkg.in/mgo.v2", true},
		{"https://godoc.org/gopkg.in/mgo.v2", true},
	}
	v := url.Values{}
	m := mngo.MongoStorage{}

	fmt.Print("Test Add Link: ")
	for _, test := range testURLs {
		v.Set("user_link", test.name)
		tf := strings.NewReader(v.Encode())
		r := httptest.NewRequest("POST", "/link/add", tf)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		s.AddLink(w, r)
		if w.Code != http.StatusFound {
			t.Errorf("Home page didn't return %v", http.StatusFound)
		}
		i := 3
		b := m.CheckMultiple(test.name, i)
		if b == true {
			t.Errorf("MongoDB has multiple insertions of %s", test.name)
		}
		fmt.Print(". ")
	}
	fmt.Println()
}

func TestRedirect(t *testing.T) {
	fmt.Print("Test Link Solver: ")
	for _, test := range testCases {
		link := test.name
		path := "/r/"
		n, _ := persist.Hash(link, path)
		r := httptest.NewRequest("GET", n, nil)
		r.Header.Set("Content-Type", "text/html")
		r.Header.Add("Accept", "text/html")
		r.Header.Set("Accept", "application/xhtml+xml")
		w := httptest.NewRecorder()

		m := mux.NewRouter()
		m.HandleFunc("/r/{id}", s.Redirect)
		m.ServeHTTP(w, r)
		if w.Code != http.StatusFound {
			fmt.Printf("\nLink %s could not be solved by app\n", link)
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
		n, dbHash := persist.Hash(link, path)
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
