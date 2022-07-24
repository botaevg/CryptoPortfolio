package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	ServerAddress string
	BaseURL       string
	DataBaseDSN   string
	Salt          string
	SecretKey     string
	CoinKey       string
}

func GetConfig() (Config, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return Config{}, err
	}
	cfg := Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		BaseURL:       os.Getenv("BASE_URL"),
		DataBaseDSN:   os.Getenv("DATABASE_DSN"),
		Salt:          os.Getenv("SALT"),
		SecretKey:     os.Getenv("SECRETKEY"),
		CoinKey:       os.Getenv("COINKEY"),
	}
	return cfg, nil
}
