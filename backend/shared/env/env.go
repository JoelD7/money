package env

import (
	"os"
	"strconv"
)

func GetString(varName, defaultValue string) string {
	val, _ := os.LookupEnv(varName)

	if val == "" {
		return defaultValue
	}

	return val
}

func GetInt(varName string, defaultValue int) int {
	val, _ := os.LookupEnv(varName)

	if val == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return intValue
}
