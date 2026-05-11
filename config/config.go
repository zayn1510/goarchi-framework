package config

import (
	"os"
	"strconv"
	"time"
)

var (
	JWTSecretKey = []byte(getEnvString("JWT_SECRET_KEY", "supersecret"))
	JWTExpired   = time.Hour * time.Duration(getEnvInt("JWT_EXPIRED_TOKEN", 5))
)

func getEnvString(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return val
}

func getEnvInt(key string, defaultValue int) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	result, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return result
}