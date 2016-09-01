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
	link := "http://localhost:8080/"
	v := url.Values{}
	v.Add("user_link", link)
	tf := strings.NewReader(v.Encode())
	r := httptest.NewRequest("POST", "/link/add", tf)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	AddLink(w, r)
	if w.Code != http.StatusFound {
		t.Errorf("Home page didn't return %v", http.StatusFound)
	}
	session, err := mgo.Dial("localhost")
	if err != nil {
		t.Errorf("Could not start session in MongoDB using localhost.")
	}
	defer session.Close()
	dbData := lines{}
	err = session.DB("tsuru").C("links").Find(bson.M{"link": link}).One(&dbData)
	if err != nil {
		t.Errorf("Could not find in MongoDB the link: %s", link)
	}
	if dbData.Link != link {
		t.Errorf("Link founded was not equal to %s", link)
	}

	var testURLs = []struct {
		name   string
		expect bool
	}{
		{"", false},
		{"notalink", false},
		{"notavalidurl.com", false},
		{"http://science.nasa.gov/", true},
		{"multiple.dots.not.valid.url", false},
		{"https://godoc.org/gopkg.in/mgo.v2", true},
		{"https://godoc.org/gopkg.in/mgo.v2", true},
	}

	// tests different URLs DB insertion
	for _, test := range testURLs {
		v.Set("user_link", test.name)
		tf = strings.NewReader(v.Encode())
		r = httptest.NewRequest("POST", "/link/add", tf)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w = httptest.NewRecorder()

		fmt.Println(v)
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
	}

	// tests the number of elements returned per query
	for _, test := range testURLs {
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
	}
}

func TestRemoveLink(t *testing.T) {

}
