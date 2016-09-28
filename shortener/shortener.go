package shortener

import (
	"crypto/md5"
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
	return fmt.Sprintf("%x", md5.Sum([]byte(u.String())))
}

func validateURL(u *url.URL) error {
	valid := govalidator.IsRequestURL(u.String())
	if valid == false {
		return fmt.Errorf("%v is a invalid url", u.String()) // probably would want a bit more informational error?
	}
	return nil
}
