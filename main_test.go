package main

import (
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

	// testing different URLs to prevent adding wrong info on DB
	// tests if the number of URLs is right ??? -to do
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

	for _, test := range testURLs {
		v := url.Values{}
		v.Add("user_link", test.name)
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

		dbData := lines{}
		var exp bool

		err = session.DB("tsuru").C("links").Find(bson.M{"link": test.name}).One(&dbData)
		if err != nil {
			if exp != test.expect {
				t.Errorf("Got a %t result, instead of %t while trying to query %s", test.expect, exp, test.name)
			}
		}
	}
	// proximos testes:
	// tentar adicionar um link diversas vezes no banco, p verificar se funciona como esperado
	// verificar parseador de url com diversas URLs validas e invalidas

}

func TestRemoveLink(t *testing.T) {

}
