package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	legalChars    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	extendedChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+[]{}|:,.<>?~"
)

// GenerateRandomString will return a truly randomize string
//
// isLegalOnly will use only legal if true will only used legal alphanumeric to generate string
func GenerateRandomString(length uint, isLegalOnly bool) string {
	if length == 0 {
		return ""
	}

	charset := legalChars
	if !isLegalOnly {
		charset = extendedChars
	}
	l := len(charset)

	charsetLen := big.NewInt(int64(len(charset)))
	var sb strings.Builder
	sb.Grow(int(length))

	for i := 0; i < int(length); i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			sb.WriteByte(charset[i%l])
			continue
		}
		sb.WriteByte(charset[randomIndex.Int64()])
	}

	return sb.String()
}
