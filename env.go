package datapace

import "os"

// Env reads specified environment variable. If no value has been found,
// fallback is returned.
func Env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
