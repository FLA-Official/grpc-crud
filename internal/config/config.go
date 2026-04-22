package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        int
	DBName        string
	JWTSecretKey  string
}

var configuration *Config

func loadConfig() {
	// Read .env file and populate environment variables.
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load the env variables", err)
		os.Exit(1)
	}

	dbusername := os.Getenv("DB_USER")
	if dbusername == "" {
		fmt.Println("Database user name is required")
		os.Exit(1)
	}

	dbpassword := os.Getenv("DB_PASSWORD")
	if dbpassword == "" {
		fmt.Println("Database user password is required")
		os.Exit(1)
	}

	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		fmt.Println("Database host is required")
		os.Exit(1)
	}

	dbport := os.Getenv("DB_PORT")
	if dbport == "" {
		fmt.Println("Database Port is required")
		os.Exit(1)
	}

	DBport, err := strconv.Atoi(dbport)
	if err != nil {
		fmt.Println("PORT must be number in env")
		os.Exit(1)
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		fmt.Println("Database Name is required")
		os.Exit(1)
	}

	jwtsecret := os.Getenv("JWT_SECRET_KEY")
	if jwtsecret == "" {
		fmt.Println("JWT Secret Key is required")
		os.Exit(1)
	}

	configuration = &Config{
		DBUser:       dbusername,
		DBPassword:   dbpassword,
		DBHost:       dbhost,
		DBPort:       DBport,
		DBName:       dbname,
		JWTSecretKey: jwtsecret,
	}
}

// GetConfig returns the loaded configuration. It loads environment variables once.
func GetConfig() *Config {
	if configuration == nil {
		loadConfig()
	}
	return configuration
}
