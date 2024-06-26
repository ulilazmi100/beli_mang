package configs

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	//"github.com/joho/godotenv"
)

type Config struct {
	DbName     string
	DbPort     string
	DbHost     string
	DbUsername string
	DbPassword string

	APPPort string
	ENV     string

	JWTSecret  string
	BcryptSalt int

	//s3 to upload, all uploaded files will available just for only a day
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_S3_BUCKET_NAME    string
	AWS_REGION            string
}

var (
	config     Config
	configOnce sync.Once
)

func LoadConfig() (Config, error) {
	var err error

	configOnce.Do(func() {
		// Load the .env file before reading the environment variables.
		// err = godotenv.Load(".env")
		// if err != nil {
		// 	fmt.Printf("No .env file found or error loading it: %v\n", err)
		// }

		config = Config{
			DbName:     GetEnvOrDefault("DB_NAME", "default_db_name"),
			DbHost:     GetEnvOrDefault("DB_HOST", "localhost"),
			DbPort:     GetEnvOrDefault("DB_PORT", "5432"),
			DbUsername: GetEnvOrDefault("DB_USERNAME", "user"),
			DbPassword: GetEnvOrDefault("DB_PASSWORD", "password"),

			APPPort: GetEnvOrDefault("APP_PORT", "8080"),
			ENV:     GetEnvOrDefault("ENV", "development"),

			JWTSecret: GetEnvOrDefault("JWT_SECRET", "sec"),

			AWS_ACCESS_KEY_ID:     GetEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
			AWS_SECRET_ACCESS_KEY: GetEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
			AWS_S3_BUCKET_NAME:    GetEnvOrDefault("AWS_S3_BUCKET_NAME", "projectsprint-bucket-public-read"),
			AWS_REGION:            GetEnvOrDefault("AWS_REGION", "ap-southeast-1"),
		}

		config.BcryptSalt, err = strconv.Atoi(GetEnvOrDefault("BCRYPT_SALT", "10"))
		if err != nil {
			err = fmt.Errorf("failed to convert BCRYPT_SALT to int: %v", err)
		}
	})

	return config, err
}

// If the variable is empty, the default value is returned.
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
