package utils

import (
	"fmt"
)

// use os package to get the env variable which is already set
func Connect(key string) string {

	// return the env variable using os package
	fmt.Println(key)
	return "ok"
}
