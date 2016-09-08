package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestHome(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	Home(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
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
	session, err := mgo.Dial("localhost")
	defer session.Close()
	v := url.Values{}

	fmt.Print("Test Add Link: ")
	// tests different URLs DB insertion
	for _, test := range testURLs {
		v.Set("user_link", test.name)
		tf := strings.NewReader(v.Encode())
		r := httptest.NewRequest("POST", "/link/add", tf)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		AddLink(w, r)
		if w.Code != http.StatusFound {
			t.Errorf("Home page didn't return %v", http.StatusFound)
		}
		dbData := lines{}
		var exp bool

		// if error not nil, It should expect `exp` & `expect` as false
		err = session.DB("tsuru").C("links").Find(bson.M{"link": test.name}).One(&dbData)
		if err != nil {
			if exp != test.expect {
				t.Errorf("Got a %t result, instead of %t while trying to query %s", test.expect, exp, test.name)
			}
		}

		// tests the number of elements returned per query
		dbNum := []lines{}
		err = session.DB("tsuru").C("links").Find(bson.M{"link": test.name}).All(&dbNum)
		if test.expect == true && err != nil {
			t.Errorf("Expected to find %s, instead MongoDB status is `%s`", test.name, err)
		} else {
			checkError(err)
		}
		if len(dbNum) > 1 {
			t.Errorf("MongoDB has multiple insertions of %s", test.name)
		}
		fmt.Print(".")
	}
	fmt.Println()
}

func TestRemove(t *testing.T) {
	var testRmv = []struct {
		name string
	}{
		{"http://localhost:8080/"},
		{"http://science.nasa.gov/"},
		{"https://godoc.org/gopkg.in/mgo.v2"},
	}
	fmt.Print("Test Removing Links: ")
	for _, test := range testRmv {
		h := md5.New()
		io.WriteString(h, test.name)
		hash := string(h.Sum(nil))
		dbHash := fmt.Sprintf("%x", hash)
		n := fmt.Sprintf("/link/remove/%s", dbHash)

		r := httptest.NewRequest("GET", n, nil)
		r.Header.Set("Content-Type", "text/html")
		r.Header.Add("Accept", "text/html")
		r.Header.Set("Accept", "application/xhtml+xml")
		w := httptest.NewRecorder()

		m := mux.NewRouter()
		m.HandleFunc("/link/remove/{id}", RemoveLink)
		m.ServeHTTP(w, r)

		session, err := mgo.Dial("localhost")
		defer session.Close()
		dbData := lines{}

		err = session.DB("tsuru").C("links").Find(bson.M{"hash": dbHash}).One(&dbData)
		if err == nil {
			t.Errorf("%s not expected on Mongo", dbHash)
		}
		fmt.Print("x")
	}
	fmt.Println()
}
