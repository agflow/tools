package config

import "os"

// LoadEnv read the value of an environment variable
func LoadEnv(varName, defaultValue string) string {
	value, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}
	return value
}
