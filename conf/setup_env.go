package conf

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Config struct {
	Port      string
	DB        *gorm.DB
	JWTSecret string
	JWTTTL    int 
}

func NewEnvConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No Env file found, using system environment variables")
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "password")
	dbName := getEnv("DB_NAME", "company_profile_db")

	db := SetupDatabaseConnection(dbHost, dbPort, dbUser, dbPass, dbName)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	jwtSecret := getEnv("JWT_SECRET", "halow")
	jwtTTLStr := getEnv("JWT_TTL", "60")
	jwtTTL, err := strconv.Atoi(jwtTTLStr)
	if err != nil {
		jwtTTL = 60
	}

	return &Config{
		Port:      getEnv("APP_PORT", "8000"),
		DB:        db,
		JWTSecret: jwtSecret,
		JWTTTL:    jwtTTL,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
