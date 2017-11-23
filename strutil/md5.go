package strutil

import (
	"crypto/sha1"
	"encoding/hex"
)

//Sha1 gen hex sha1 of content
func Sha1(content string) string {
	hash := sha1.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}
