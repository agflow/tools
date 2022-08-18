package agerr

import "fmt"

// Wrap wraps an error with `wrapmsg`
func Wrap(wrapmsg string, err error) error {
	if err != nil {
		return fmt.Errorf(wrapmsg, err)
	}
	return nil
}
