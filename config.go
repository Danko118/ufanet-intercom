package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	MQTTBroker string
	AppPort    string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // загружает .env, игнорирует ошибку если файл не найден

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		MQTTBroker: os.Getenv("MQTT_BROKER"),
		AppPort:    os.Getenv("APP_PORT"),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		logger.WithFields(logrus.Fields{
			"state":   "Init",
			"status":  "Error",
			"service": "CFG",
		}).Fatal("Не удалось прочитать env переменные")
	}

	return cfg
}
