package utils

import (
	"regexp"
	"strconv"
)

// Utilities
func Extract(re *regexp.Regexp, text string) string {
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func UnquoteUnicode(str string) string {
	if str == "" {
		return ""
	}
	s, err := strconv.Unquote(`"` + str + `"`)
	if err != nil {
		return str
	}
	return s
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
