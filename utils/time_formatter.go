package utils

import "strings"

func CutTimeString(s string) string {
	before, _, _ := strings.Cut(s, "m")
	before = strings.Replace(before, "h", " ч. ", 1)
	return before + " м."
}
