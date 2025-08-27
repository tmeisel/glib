package strutils

import (
	cryptoRand "crypto/rand"
	"math/big"
	mathRand "math/rand"
	"slices"
	"strings"
)

// Ptr returns a pointer to s
func Ptr(s string) *string {
	return &s
}

// InSlice returns true, if the given string needle is contained in the
// slice of strings haystack
func InSlice(needle string, haystack []string) bool {
	return slices.Contains(haystack, needle)
}

// InSliceIgnoreCase returns true, if the given string needle is in the
// slice of strings haystack, ignoring the case of the comparable
func InSliceIgnoreCase(needle string, haystack []string) bool {
	for _, elem := range haystack {
		if strings.ToLower(needle) == strings.ToLower(elem) {
			return true
		}
	}

	return false
}

// SubString returns a portion of a given string. If length is 0,
// everything from pos until the end of the given string will be returned
func SubString(s string, pos, length int) string {
	runes := []rune(s)

	total := len(runes)

	if pos > total {
		return ""
	}

	if length == 0 {
		return string(runes[pos:])
	}

	if length < 0 {
		pos = total + length
		length = pos + (length * -1)

		return string(runes[pos:length])
	}

	l := pos + length
	if l > total {
		l = total
	}

	return string(runes[pos:l])
}

// Shuffle shuffles the given string s' characters
func Shuffle(s string) string {
	output := make([]byte, len(s))
	perm := mathRand.Perm(len(s) - 1)
	for idx, rnd := range perm {
		output[idx] = s[rnd]
	}

	return string(output)
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

// Random returns a random string generated using crypto/cryptoRand.
// The alphabet can be one or more strings, including the ones
// provided by this pkg like AlphabetNumbers.
func Random(length int, alphabet ...string) (string, error) {
	chars := strings.Join(alphabet, "")
	outputChars := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}

		outputChars[i] = chars[num.Int64()]
	}

	return Shuffle(string(outputChars)), nil
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
