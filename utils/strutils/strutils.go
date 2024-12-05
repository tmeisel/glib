package strutils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// Ptr returns a pointer to s
func Ptr(s string) *string {
	return &s
}

// SubString returns a portion of a given string. If length is 0,
// everything from pos until the end of the given string will be returned
func SubString(s string, pos, length int) string {
	runes := []rune(s)

	if pos > len(runes) {
		return ""
	}

	if length == 0 {
		return string(runes[pos:])
	}

	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}

const (
	// AlphabetUCChars Upper case chars
	AlphabetUCChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// AlphabetLCChars // Lower case chars
	AlphabetLCChars      = "abcdefghijklmnopqrstuvwxyz"
	AlphabetNumbers      = "1234567890"
	AlphabetSpecialChars = ",;.:-_!§$%&/()=?*+#"

	// Alphabet is a combination of upper and lower case chars, numbers and special chars
	Alphabet = AlphabetLCChars +
		AlphabetUCChars +
		AlphabetNumbers +
		AlphabetSpecialChars

	// AlphabetReadable is an alphanumeric alphabet that excludes chars that can
	// be confused, depending on the used font. E.g. 0 and O or l and I
	AlphabetReadable = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKKLMNPQRSTUVWXYZ2345689"
)

// Random returns a random string generated using crypto/rand.
// The alphabet can be one or more strings, including the ones
// provided by this pkg like AlphabetNumbers.
func Random(length int, alphabet ...string) (string, error) {
	chars := strings.Join(alphabet, "")
	output := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}

		output[i] = chars[num.Int64()]
	}

	return string(output), nil
}

// MustRandom returns a random string. For a detailed description,
// see Random
func MustRandom(length int, alphabet ...string) string {
	str, err := Random(length, alphabet...)
	if err != nil {
		panic(err)
	}

	return str
}
