package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   string
	JWTSecretKey string
	BcryptCost   int
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	DatabaseURL  string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	bcryptCost := 12
	bcryptCostStr := os.Getenv("BCRYPT_COST")

	if bcryptCostStr != "" {
		parsedCost, err := strconv.Atoi(bcryptCostStr)
		if err != nil {
			log.Printf("Invalid BCRYPT_COST value '%s', using default 12", bcryptCostStr)
		} else {
			bcryptCost = parsedCost
		}
	}

	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePass := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")

	databaseURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		databaseHost, databasePort, databaseUser, databasePass, databaseName)

	return &Config{
		ServerPort:   os.Getenv("SERVER_PORT"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
		BcryptCost:   bcryptCost,
		DatabaseHost: databaseHost,
		DatabasePort: databasePort,
		DatabaseUser: databaseUser,
		DatabasePass: databasePass,
		DatabaseName: databaseName,
		DatabaseURL:  databaseURL,
	}
}
