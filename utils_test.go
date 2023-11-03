package errorst

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func hideLineNumbers(s string) string {
	digits := regexp.MustCompile(`\d+`)
	return digits.ReplaceAllString(s, "##")
}

func pathWithPrefix(path string) string {
	// get the path to the current file
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.ToSlash(dir) + "/" + path
}

func formatLoc(file, fun string) string {
	return fmt.Sprintf(" --- at %v:## (%v) ---", file, fun)
}
