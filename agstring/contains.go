package agstring

import "github.com/thoas/go-funk"

// ContainsAny checks if source slice contains any of given strings
func ContainsAny(src []string, ss ...string) bool {
	for _, s := range ss {
		if funk.ContainsString(src, s) {
			return true
		}
	}
	return false
}

// ContainsAll checks if source slice contains all the given strings
func ContainsAll(src []string, ss ...string) bool {
	for _, s := range ss {
		if !funk.ContainsString(src, s) {
			return false
		}
	}
	return true
}
