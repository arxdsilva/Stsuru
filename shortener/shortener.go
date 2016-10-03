package shortener

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/asaskevich/govalidator"
)

// Shorten does the hard work about making your url small
func Shorten(u *url.URL, customHost string) (*url.URL, error) {
	err := validateURL(u)
	if err != nil {
		return nil, err
	}
	hash := hashGenerator(u)
	switch customHost {
	case "":
		return &url.URL{
			Scheme: "https",
			Host:   u.Host,
			Path:   hash,
		}, nil
	default:
		return &url.URL{
			Scheme: "https",
			Host:   customHost,
			Path:   hash,
		}, nil
	}
}

func hashGenerator(u *url.URL) string {
	hasher := md5.New()
	hasher.Write([]byte(u.String()))
	return hex.EncodeToString(hasher.Sum(nil))
}

func validateURL(u *url.URL) error {
	v := govalidator.IsURL(u.String())
	valid := govalidator.IsRequestURL(u.String())
	if !valid || !v {
		return fmt.Errorf("%v is a invalid url", u.String())
	}
	return nil
}

func tokenGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
