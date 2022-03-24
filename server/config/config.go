package config

import (
	"log"
	"os"
	"path/filepath"

	"pfserver/utils"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

const (
	BEAMAAN_FEES = 6 // 6%
)

func init() {
	// in prod we won't be using the .env, maybe use Vault
	if !IsProd() {
		err := godotenv.Load(GetEnvPath())

		if err != nil {
			log.Fatal("Loading Env Error: ", err)
		}
	}
}

func GetEnvPath() string {
	return filepath.Join(utils.RootPath(), ".env")

}

func Cors() cors.Options {
	return cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://secure-global.paytabs.com",
			"https://beamaanapp.vercel.app",
		},
		AllowCredentials: true,
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

// check if we are in testing mode (go run main.go testing=true)
func IsTestMode() bool {
	return os.Args[len(os.Args)-1] == "testing=true"
}

func IsProd() bool {
	return GetEnv("NODE_ENV") == "production"
}
