package errorst_test

import (
	"regexp"
)

func hideLineNumbers(s string) string {
	digits := regexp.MustCompile(`\d+`)
	return digits.ReplaceAllString(s, "##")
}
