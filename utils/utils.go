package utils

import (
	"crypto/sha1"
	"io"
	"fmt"
)

//EncryptPassword makes password encryption with SHA1 algorithm
func EncryptPassword(password string) string {
	h := sha1.New()
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}