package strutils

func SubString(s string, pos, length int) string {
	runes := []rune(s)

	if pos > len(runes) {
		return ""
	}

	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}

	return string(runes[pos:l])
}
