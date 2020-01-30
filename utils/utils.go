package utils

import (
	"os"
)

// use os package to get the env variable which is already set
func EnvVariable(key string) string {

	// return the env variable using os package
	return os.Getenv(key)
}
