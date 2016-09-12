package mngo

import (
	"testing"
)

var test = []struct {
	link  string
	hash  string
	isURL bool
}{
	{"abcdef", "", false},
	{"www.notvalidurl.netscape", "", false},
	{"http://www.gorillatoolkit.org/pkg/mux", "", true},
}

func TestInsert(t *testing.T) {
	Insert(link)
}

func TestFindLink(t *testing.T) {
	FindLink(s)
}

func TestFindHash(t *testing.T) {
	FindHash(s)
}

func TestGetAll(t *testing.T) {
	GetAll()
}

func TestHash(t *testing.T) {
	Hash(link, path)
}
