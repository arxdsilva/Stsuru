package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	_ "github.com/virtao/GoEndian"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	KeySize   = 32
	NonceSize = 24
)

// GenerateKey creates a new random secret key.
func GenerateKey() (*[KeySize]byte, error) {
	key := new([KeySize]byte)
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GenerateNonce creates a new random nonce.
func GenerateNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

var (
	ErrEncrypt = errors.New("secret: encryption failed")
	ErrDecrypt = errors.New("secret: decryption failed")
)

func Encrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, ErrEncrypt
	}

	out := make([]byte, len(nonce))
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, key)
	return out, nil
}

func Decrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	if len(message) < (NonceSize + secretbox.Overhead) {
		return nil, ErrDecrypt
	}

	var nonce [NonceSize]byte
	copy(nonce[:], message[:NonceSize])
	out, ok := secretbox.Open(nil, message[NonceSize:], &nonce, key)
	if !ok {
		return nil, ErrDecrypt
	}

	return out, nil
}

func main() {

	buff := make([]byte, 8)
	key, _ := GenerateKey()

	var length int

	roundTrip(key, buff, 1)
	roundTrip(key, buff, 1000)
	roundTrip(key, buff, 2000000)
	roundTrip(key, buff, 4294967295)
	roundTrip(key, buff, 10000000000000)
	roundTrip(key, buff, 50000000000000000)
	roundTrip(key, buff, 18446744073709551615)

	fmt.Println(length)
}

func roundTrip(key *[32]byte, buff []byte, value uint64) {
	fmt.Println("processing", value)
	//network byte order is best for portability
	binary.BigEndian.PutUint64(buff, value)
	secret, _ := Encrypt(key, buff)
	encodedSecret := base64.URLEncoding.EncodeToString(secret)
	fmt.Println("encoded:", encodedSecret, len(encodedSecret), "bytes")
	clear, _ := Decrypt(key, secret)
	fmt.Println("clear:", clear, len(clear), "bytes")
}
