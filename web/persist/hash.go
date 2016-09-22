package persist

import (
	"crypto/md5"
	"fmt"
	"io"
)

// Hash creates & returns a link with the hashed URL and the URL hash
func Hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%s%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}
