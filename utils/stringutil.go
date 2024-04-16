package utils

func EndsWith(s string, end string) bool {
	return s[:1] == end
}

func StartsWith(s string, end string) bool {
	return s[1:] == end
}
