package common

import (
	"os"
)

func GetEnvWithString(name string, defaultValue string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
