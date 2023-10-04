package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type EnvType string

const (
	Dev        EnvType = `development`
	Test               = `test`
	Production         = `production`
)

type EnvConfig struct {
	Type EnvType
}

var envOnce sync.Once
var envConfig EnvConfig

func envInstance() EnvConfig {
	envOnce.Do(func() {
		_ = godotenv.Overload(`.env`)
		envName, exists := os.LookupEnv("ENV")
		if !exists {
			panic("ENV variable is not set. possible values are: development, test, production")
		}
		et, err := parseEnvName(envName)
		envConfig = EnvConfig{Type: et}
		if err != nil {
			panic(err)
		}
		return
	})
	return envConfig
}

func (ec EnvConfig) getEnvValue(key string, defValue string) string {
	envValue, exists := os.LookupEnv(key)
	if !exists {
		return defValue
	}
	return envValue
}
func parseEnvName(envName string) (EnvType, error) {
	switch EnvType(strings.ToLower(envName)) {
	case Dev:
		return Dev, nil
	case Test:
		return Test, nil
	case Production:
		return Production, nil
	default:
		return "", errors.New(fmt.Sprintf("invalid env name:%s. possible values are: development, test, production",
			envName))
	}
}

func GetEnvValueOrPanic(key string) string {
	value := GetEnvValue(key, "")
	if value == "" {
		log.Fatal(fmt.Sprintf("%s env variable is undefined", key))
	}
	return value
}

func GetEnvValue(key, defValue string) string {
	return envInstance().getEnvValue(key, defValue)
}

func GetEnvValueBool(key string, defValue bool) bool {
	val := GetEnvValue(key, "")
	if val == "" {
		return defValue
	} else if strings.ToLower(strings.Trim(val, " ")) == "true" {
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
	return envInstance().Type == Production
}

func IsDevelopment() bool {
	return envInstance().Type == Dev
}

func IsTesting() bool {
	return envInstance().Type == Test
}
