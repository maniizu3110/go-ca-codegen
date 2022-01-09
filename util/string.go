package util

import "strings"

func SamePrefix(a, b string) string {
	for i := range a {
		if !strings.HasPrefix(b, a[:i+1]) {
			return a[:i]
		}
	}
	return a
}
func SameSuffix(a, b string) string {
	for i := range a {
		if !strings.HasSuffix(b, a[len(a)-i-1:]) {
			return a[len(a)-i:]
		}
	}
	return a
}
