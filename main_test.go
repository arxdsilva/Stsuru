package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

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
}

func TestRemoveLink(t *testing.T) {

}
