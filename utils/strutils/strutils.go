package strutils

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
