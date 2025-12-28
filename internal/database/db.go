package database

import (
	"os"

	"github.com/Sandstorm831/021_trade/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	logrus.Info("Successfully connected to DB")

}

func MigrateToDB() {
	logrus.Info("Starting migrations")
	DB.AutoMigrate(&models.User{}, &models.Stock{}, &models.StockPrice{}, &models.Reward{}, models.LedgerEntry{})
	logrus.Info("Applied migration successuflly")
}
