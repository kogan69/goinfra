package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func InitEnvFromFile(envFilePath string) (err error) {
	err = godotenv.Overload(envFilePath)
	if err != nil {
		return
	}
	return err
}

func GetEnvValueOrPanic(key string) string {
	value := GetEnvValue(key, "")
	if value == "" {
		log.Fatal(fmt.Sprintf("%s env variable is undefined", key))
	}
	return value
}
func GetEnvValue(key, defValue string) string {
	envValue, exists := os.LookupEnv(key)
	if !exists {
		return defValue
	}
	return envValue
}

func GetEnvValueBool(key string, defValue bool) bool {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defValue
	}
	if strings.ToLower(strings.Trim(val, " ")) == "true" {
		return true
	} else {
		return false
	}
}

func GetEnvValueInt(key string, defValue int) int {
	sv := GetEnvValue(key, "")
	if sv == "" {
		return defValue
	} else {
		v, err := strconv.Atoi(sv)
		if err != nil {
			return defValue
		} else {
			return v
		}
	}

}

func GetEnvValueFloat64(key string, defValue float64) float64 {
	sv := GetEnvValue(key, "")
	if sv == "" {
		return defValue
	} else {
		v, err := strconv.ParseFloat(sv, 64)
		if err != nil {
			return defValue
		} else {
			return v
		}
	}

}

func IsProduction() bool {
	return GetEnvValue("ENV", "development") == "production"
}

func IsDevelopment() bool {
	return GetEnvValue("ENV", "development") == "development"
}

func IsTesting() bool {
	return GetEnvValue("ENV", "development") == "test"
}
