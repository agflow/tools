package config

import (
	"os"

	"github.com/agflow/tools/agnumber"
)

// LoadEnv loads the value of an environment variable
func LoadEnv(varName, defaultValue string) string {
	value, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}
	return value
}

// LoadIntEnv loads the value of an environment variable into a int variable
func LoadIntEnv(varName string, defaultValue int64) int64 {
	value, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}
	valueInt, err := agnumber.Atoi64(value)
	if err != nil {
		return defaultValue
	}
	return valueInt
}
