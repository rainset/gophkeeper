package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func Md5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func Sha256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
