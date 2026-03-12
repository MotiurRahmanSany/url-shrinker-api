package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var serviceConfig *Config

type DBConfig struct {
	Host          string
	Port          int
	User          string
	Password      string
	Name          string
	EnableSSLMode bool
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type Config struct {
	HttpPort     int
	Version      string
	AppName      string
	JwtSecretKey string
	Db           *DBConfig
	Redis        *RedisConfig
}

func loadConfig() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file, " + err.Error())
	}

	version := os.Getenv("VERSION")

	if version == "" {
		fmt.Println("VERSION not set in .env file, using default value 1.0.0")
	}

	httpPort := os.Getenv("HTTP_PORT")

	if httpPort == "" {
		fmt.Println("HTTP_PORT not set in .env file, using default value 8080")
	}

	httpPortInt, err := strconv.ParseInt(httpPort, 10, 64)

	if err != nil {
		fmt.Printf("Error parsing HTTP_PORT: %v, using default value 8080\n", err)
	}

	serviceName := os.Getenv("SERVICE_NAME")

	if serviceName == "" {
		fmt.Println("SERVICE_NAME not set in .env file, using default value Student Management System")
	}

	jwtSEcretKey := os.Getenv("JWT_SECRET_KEY")

	if jwtSEcretKey == "" {
		fmt.Println("JWT_SECRET_KEY not set in .env file, using default value very_strong_secret_key123")
	}

	dbHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbEnableSSLMode := os.Getenv("DB_SSL_MODE")

	if dbHost == "" || DBPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbEnableSSLMode == "" {
		fmt.Println("One or more database configuration values are missing in .env file, using default values, Please set these values first")

		os.Exit(1)
	}

	dbPortInt, err := strconv.ParseInt(DBPort, 10, 64)

	if err != nil {
		fmt.Printf("Error parsing DB_PORT: %v, using default value 5432\n", err)
		os.Exit(1)
	}

	dbEnableSSLModeBool, err := strconv.ParseBool(dbEnableSSLMode)

	if err != nil {
		fmt.Printf("Error parsing DB_SSL_MODE: %v, using default value false\n", err)
		os.Exit(1)
	}

	dbConfig := &DBConfig{
		Host:          dbHost,
		Port:          int(dbPortInt),
		User:          dbUser,
		Password:      dbPassword,
		Name:          dbName,
		EnableSSLMode: dbEnableSSLModeBool,
	}

	// Redis config
	redisHost := os.Getenv("REDIS_HOST")
	redisPortStr := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	if redisHost == "" || redisPortStr == "" || redisDBStr == "" {
		fmt.Println("One or more Redis values missing in .env")
		os.Exit(1)
	}

	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		fmt.Printf("Error parsing REDIS_PORT: %v, using default 6379\n", err)
		redisPort = 6379
	}

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		fmt.Printf("Error parsing REDIS_DB: %v, using default 0\n", err)
		redisDB = 0
	}

	redisConfig := &RedisConfig{
		Host:     redisHost,
		Port:     redisPort,
		Password: redisPassword,
		DB:       redisDB,
	}

	serviceConfig = &Config{
		HttpPort:     int(httpPortInt),
		Version:      version,
		AppName:      serviceName,
		JwtSecretKey: jwtSEcretKey,
		Db:           dbConfig,
		Redis:        redisConfig,
	}

}

//  returns the loaded config singleton
func GetConfig() *Config {
	if serviceConfig == nil {
		loadConfig()
	}

	return serviceConfig
}
