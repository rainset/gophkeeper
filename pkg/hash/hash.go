package hash

import (
	"crypto/md5"
	"crypto/sha256"
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
