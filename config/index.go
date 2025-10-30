package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type envConfig struct {
	AppPort string
	DbPath string
	GoogleClientID string
	GoogleClientSecret string
	GoogleRedirectURL string
	JWTSaltKey string
	FEOriginURL string
}

func (e *envConfig) LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("ENV file not loaded")
	}
	e.AppPort = loadString("APP_PORT", ":8080")
	e.DbPath = loadString("DB_PATH", "")
	e.GoogleClientID = loadString("GOOGLE_CLIENT_ID", "")
	e.GoogleClientSecret = loadString("GOOGLE_CLIENT_SECRET", "")
	e.GoogleRedirectURL = loadString("GOOGLE_REDIRECT_URL", "")
	e.JWTSaltKey = loadString("JWT_SALT_KEY", "test_salt")
	e.FEOriginURL = loadString("FE_ORIGIN_URL", "")
}

var Config envConfig

func init() {
	Config.LoadConfig()
}

func loadString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		log.Panic("APP_PORT is not loaded")
		return fallback
	}
	return val
}