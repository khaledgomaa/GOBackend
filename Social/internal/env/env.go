// Equivalent to IConfiguration in .NET
// Provides helpers to read environment variables with fallbacks.
package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}

func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return result
}
