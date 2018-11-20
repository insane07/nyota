//
// Provides helper funcs to get environ variable values of different types.
package sysutils

import (
	"os"
	"strconv"
	"time"
)

func Getenv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return defaultValue
}

func GetenvInt(name string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err == nil {
		return value
	}
	return defaultValue
}

func GetenvBool(name string, defaultValue bool) bool {
	value, err := strconv.ParseBool(os.Getenv(name))
	if err == nil {
		return value
	}
	return defaultValue
}

func GetenvTime(name string, unit time.Duration, defaultValue int) time.Duration {
	v := GetenvInt(name, defaultValue)
	return time.Duration(v) * unit
}
