package env

import "os"

func GetString(varName, defaultValue string) string {
	val, _ := os.LookupEnv(varName)

	if val == "" {
		return defaultValue
	}

	return val
}
