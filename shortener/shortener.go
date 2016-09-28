package shortener

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/asaskevich/govalidator"
)

// NewLink ...
func NewLink(link, path string) (string, string, error) {
	v := validateURL(link)
	if !v {
		return "", "", fmt.Errorf("Given link: `%s` is Not a valid URL", link)
	}
	hashedPath, hashNum := hash(link, path)
	return hashNum, hashedPath, nil
}

func hash(link, path string) (string, string) {
	h := md5.New()
	io.WriteString(h, link)
	hash := string(h.Sum(nil))
	linkShort := fmt.Sprintf("%s%x", path, hash)
	dbHash := fmt.Sprintf("%x", hash)
	return linkShort, dbHash
}

func validateURL(l string) bool {
	isURL := govalidator.IsURL(l)
	validURL := govalidator.IsRequestURL(l)
	if isURL == false || validURL == false {
		return false
	}
	return true
}
