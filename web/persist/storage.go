package persist

import "github.com/arxdsilva/Stsuru/web/persist/data"

// Storage is the interface that holds how the server uses the given data.
type Storage interface {
	List() ([]Stored, error)
	Save(linkData *data.LinkData) error
	Remove(dbHash string) error
	FindHash(dbHash string) (string, error)
}
