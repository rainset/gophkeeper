package hash

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func Md5(text string) string {
	hash := md5.Sum([]byte(text))
	strHash := hex.EncodeToString(hash[:])

	return strHash
}

func Sha256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	strHash := hex.EncodeToString(h.Sum(nil))

	return strHash
}

// GenerateRandom generates slice of random bytes with specified size.
// GenerateRandom генерирует случайную последовательность байт длинной size
func GenerateRandom(size int) ([]byte, error) {
	// генерируем случайную последовательность байт
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
