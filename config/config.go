package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func LoadConfig() error {
	env := map[string]string{
		"DEV": ".dev.env",
		"":    ".dev.env",
	}

	key := strings.ToUpper(os.Getenv("APP_ENV"))
	envPath, ok := env[key]
	if !ok {
		envPath = env[""]
	}

	err := godotenv.Load(basepath + "/" + envPath)
	if err != nil {
		return err
	}

	return nil
}
