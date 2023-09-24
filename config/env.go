package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

func InitEnv(appEnv string) {
	InitEnvFromPath(appEnv, "")
}

// InitEnvFromPath /**
func InitEnvFromPath(appEnv, path string) {
	if len(path) != 0 {
		path = strings.TrimRight(path, "/") + "/"
	}
	envFiles := make([]string, 0)
	if len(appEnv) == 0 {
		envFiles = append(envFiles, path+".env", path+".env.local")
	} else {
		envFiles = append(envFiles, path+".env", path+".env.local", path+".env."+appEnv, path+".env."+appEnv+".local")
	}
	load(envFiles)
}
func GetEnvValueOrPanic(key string) string {
	value := GetEnvValue(key, "")
	if value == "" {
		log.Fatal(fmt.Sprintf("%s env variable is undefined", key))
	}
	return value
}
func GetEnvValue(key, defValue string) string {
	mutex.RLock()
	defer mutex.RUnlock()
	if envValue, exists := env[key]; exists {
		return envValue
	}
	if envValue, exists := os.LookupEnv(key); exists {
		env[key] = envValue
		return envValue
	}
	return defValue
}

func PrintOutEnv(file *os.File) {
	_, _ = fmt.Fprintf(file, "environment variables\n")
	for k, v := range env {
		_, _ = fmt.Fprintf(file, "%s=%s\n", k, v)
	}
}

func GetEnvValueBool(key string, defValue bool) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	val, exists := env[key]
	if !exists {
		if v, ok := os.LookupEnv(key); ok {
			env[key] = v
			val = v
		}
	}
	if val == "" {
		return false
	} else {
		if strings.ToLower(strings.Trim(val, " ")) == "true" {
			return true
		} else {
			return false
		}
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

func SetEnvValue(key, value string) error {
	mutex.Lock()
	defer mutex.Unlock()
	err := os.Setenv(key, value)
	if err == nil {
		env[key] = value
	}
	return err
}

func IsProduction() bool {
	return GetEnvValue("API_ENV", "development") == "production"
}

func IsDevelopment() bool {
	return GetEnvValue("API_ENV", "development") == "development"
}

func IsTesting() bool {
	return GetEnvValue("API_ENV", "development") == "test"
}

func IsDebug() bool {
	return GetEnvValue("API_ENV", "debug") == "debug"
}

var mutex = &sync.RWMutex{}
var env = map[string]string{}

// load .env files. Files will be loaded in the same order that are received.
// Redefined vars will override previously existing values.
// IE: envy.load(".env", "test_env/.env") will result in DIR=test_env
// If no arg passed, it will try to load a .env file.
func load(files []string) {
	// If no files received, load the default one ".env"
	if len(files) == 0 {
		err := godotenv.Overload()
		if err == nil {
			reload()
		}
		return
	}

	for _, file := range files {
		// Check if it exists or we can access
		if _, err := os.Stat(file); err != nil {
			// It does not exist or we can not access.
			continue
		}

		// It exists and we have permission. load it
		if err := godotenv.Overload(file); err != nil {
			continue
		}

		// Reload the env so all new changes are noticed
		reload()
	}
}

// Reload the ENV variables. Useful if
// an external ENV manager has been used
func reload() {
	mutex.Lock()
	defer mutex.Unlock()
	env = map[string]string{}
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		env[pair[0]] = os.Getenv(pair[0])
	}
}
