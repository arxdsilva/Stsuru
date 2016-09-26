package persist

// Storage ...
type Storage interface {
	List() ([]Stored, error)
	Save(link, linkShort, dbHash string) error
	Remove(dbHash string) error
	FindHash(dbHash string) (string, error)
}
