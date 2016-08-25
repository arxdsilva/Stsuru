package main

import (
	"net/http/httptest"
	"testing"
)

func requestTest(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	Home(w, req)
}
