package model

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/thoas/go-funk"

	"github.com/agflow/tools/agstring"
)

// User struct of user received from account management
type User struct {
	Data map[string]interface{} `json:"data,omitempty"`
}

func toStringSlice(ls ...interface{}) ([]string, error) {
	result := make([]string, len(ls))
	for i := 0; i < len(ls); i++ {
		var ok bool
		result[i], ok = ls[i].(string)
		if !ok {
			return nil, errors.New("can't cast user rule to string")
		}
	}
	return result, nil
}

// GetRules returns rules for User
func (u *User) GetRules() (Rules, error) {
	rules, ok := u.Data["rules"]
	if !ok {
		return nil, errors.New("can't get rules")
	}
	return toStringSlice(rules.([]interface{})...)
}

func (u *User) ContainsRules(rules ...string) bool {
	rules, err := u.GetRules()
	if err != nil {
		return false
	}
	return agstring.ContainsAll(rules, rules...)
}

// Rules is array of string
type Rules []string

// FilterByPrefix filters rules by prefix
func (rs *Rules) FilterByPrefix(s string) Rules {
	return funk.FilterString(*rs, func(rule string) bool {
		return strings.HasPrefix(rule, s)
	})
}

func (rs *Rules) Hash() (string, error) {
	h := sha512.New()
	for _, s := range *rs {
		if _, err := io.WriteString(h, s); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Contains returns true is rules contain given string
func (rs *Rules) Contains(s string) bool {
	return funk.ContainsString(*rs, s)
}
