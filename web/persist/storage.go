package persist

// Storage ...
type Storage interface {
	List() ([]Stored, error)
	Save(link, customHost, dbHash string) error
	Remove(dbHash string) error
	FindHash(dbHash string) (string, error)
}
