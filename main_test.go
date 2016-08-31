package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	Home(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}

// func TestAddLink(t *testing.T) {
// 	req := httptest.NewRequest("POST", "/", nil)
// 	w := httptest.NewRecorder()
// 	AddLink(w, req)
// 	if w.Code != http.StatusNotFound {
// 		t.Errorf("Home page didn't return %v", http.StatusOK)
// 	}
// }
