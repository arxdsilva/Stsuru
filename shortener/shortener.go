package shortener

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/asaskevich/govalidator"
)

type NewShorten struct {
	U          *url.URL
	CustomHost string
	Token      string
	NumBytes   int
}

// Shorten recieves your customHost and applies to the url to be returned, so
func (n *NewShorten) Shorten() (*url.URL, error) {
	err := validateURL(n.U)
	if err != nil {
		return nil, err
	}
	hash := switchToken(n.U, n.Token, n.NumBytes)
	return switchHost(n.U, hash, n.CustomHost)
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

func tokenGenerator(numBytes int) string {
	switch numBytes {
	case 0:
		numBytes = 4
	}
	b := make([]byte, numBytes)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func switchToken(u *url.URL, s string, n int) string {
	switch s {
	case "":
		return tokenGenerator(n)
	default:
		return hashGenerator(u)
	}
}

func switchHost(u *url.URL, hash, customHost string) (*url.URL, error) {
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
