package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
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
