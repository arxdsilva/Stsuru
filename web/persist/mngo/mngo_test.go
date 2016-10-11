package mngo

import (
	"fmt"
	"testing"

	"github.com/arxdsilva/Stsuru/web/persist/data"
)

var testCases = []struct {
	link, linkShort, hash string
	isURL                 bool
}{
	{
		"abcdef",
		"",
		"",
		false,
	},
	{
		"www.globo.com",
		"",
		"",
		false,
	},
	{
		"www.notvalidurl.netscape",
		"",
		"",
		false,
	},
	{
		"http://www.gorillatoolkit.org/pkg/mux",
		"",
		"70df8650c03c9fdfc959f04a64ecd956",
		true,
	},
	{
		"https://mail.google.com/mail/u/0/#inbox",
		"",
		"2122c5656da3d86d77c08f7af48c0268",
		true,
	},
	{
		"https://mail.google.com/mail/u/0/#inbox",
		"",
		"2122c5656da3d86d77c08f7af48c0268",
		true,
	},
	{
		"https://www.youtube.com/watch?v=grwx4OMfAn4",
		"",
		"678989a28d9b88ada6cc6678df8e6aa1",
		true,
	},
}

var s = MongoStorage{
	URL:        "localhost",
	DB:         "tsuru",
	Collection: "links",
}

func TestSave(t *testing.T) {
	fmt.Print("Test Save: ")
	for _, test := range testCases {
		if !test.isURL {
			InsertData := data.LinkData{
				Link:  test.link,
				Short: test.linkShort,
				Hash:  test.hash,
			}
			err := s.Save(&InsertData)
			if err != nil {
				t.Fatalf("unexpected save error: %v", err)
			}
			fmt.Print(".")
		}
	}
	fmt.Println()
}

func TestFindHash(t *testing.T) {
	fmt.Print("Test FindHash: ")
	for _, test := range testCases {
		link, err := s.FindHash(test.hash)
		if err != nil && !test.isURL {
			fmt.Print(".")
			continue
		}
		if test.link == link {
			fmt.Print(".")
		}
		continue
	}
	fmt.Println()
}

func TestGetAll(t *testing.T) {
	// Teste esta falhando pois estou usando o MongoStorage da linha 56, quando deveria usar persist.Storage ... (p/ usar os comportamentos de Storage)
	// Melhorar mensagens de erro nesse teste.
	fmt.Print("Test GetAll: ")
	a, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error getall: %v", err)
	}
	if len(a) == 3 {
		fmt.Print("...")
	} else {
		t.Errorf("Array bigger than expected, len == %v, expected = %v", len(a), 3)
	}
	fmt.Println()
}

func TestRemove(t *testing.T) {
	fmt.Print("Test Remove: ")
	for _, test := range testCases {
		err := s.Remove(test.hash)
		if err != nil {
			continue
		}
		fmt.Print(".")
	}
	fmt.Println()
}
