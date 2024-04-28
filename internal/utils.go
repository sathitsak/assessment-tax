package internal

import "os"

func GetEnvWithFallback(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
