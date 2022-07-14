package security

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
)

// Hash hashes `vs`
func Hash(vs interface{}) (string, error) {
	h := sha512.New()
	r, err := json.Marshal(vs)
	if err != nil {
		return "", err
	}
	h.Write(r)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
