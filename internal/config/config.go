package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func LoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		panic("Error laoding .env file")
	}
	logrus.Info("Environement Variables loading Successful")
}
