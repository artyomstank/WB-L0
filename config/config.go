package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer HTTPServer
	Postgres   Postgres
	Cache      Cache
	Kafka      Kafka
}

type HTTPServer struct {
	Host    string
	Port    int
	Timeout time.Duration
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Cache struct {
	StartupSize int
}

type Kafka struct {
	Host  string
	Port  int
	Topic string
	Group string
}

// Подгружаем .env, если есть
func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPServer: HTTPServer{
			Host:    getEnv("HTTP_HOST", "localhost"),
			Port:    getEnvAsInt("HTTP_PORT", 8081),
			Timeout: getEnvAsDuration("HTTP_TIMEOUT", 5*time.Second),
		},
		Postgres: Postgres{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5433),
			User:     getEnv("POSTGRES_USER", "wb_tech_user"),
			Password: getEnv("POSTGRES_PASSWORD", "12345678"),
			Database: getEnv("POSTGRES_DATABASE", "wb_tech_demo_service"),
		},
		Cache: Cache{
			StartupSize: getEnvAsInt("CACHE_STARTUP_SIZE", 10),
		},
		Kafka: Kafka{
			Host:  getEnv("KAFKA_HOST", "localhost"),
			Port:  getEnvAsInt("KAFKA_PORT", 9092),
			Topic: getEnv("KAFKA_TOPIC", "wb-orders"),
			Group: getEnv("KAFKA_GROUP", "wb-tech-demo-service"),
		},
	}

	return cfg
}

// Получить строку подключения к БД
func (p Postgres) GetDBConnStr() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Database,
	)
}

// вспомогательные функции подгрузки конфига
func GetLimitCache() int {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются стандартные переменные окружения")
	}
	limitStr := os.Getenv("ORDERS_LIMIT")
	limit, _ := strconv.Atoi(limitStr)
	return limit
}

func getEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr, exists := os.LookupEnv(key); exists {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
		log.Printf("не удалось преобразовать %s=%s в целое число, используется значение по умолчанию %d", key, valStr, defaultVal)
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr, exists := os.LookupEnv(key); exists {
		if val, err := time.ParseDuration(valStr); err == nil {
			return val
		}
		log.Printf("не удалось преобразовать %s=%s в значение длительности, используется значение по умолчанию %s", key, valStr, defaultVal)
	}
	return defaultVal
}
