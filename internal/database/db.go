package database

import (
	"fmt"
	"os"

	"github.com/Sandstorm831/021_trade/internal/models"
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
	fmt.Println("Connected to DB")

}

func MigrateToDB() {
	DB.AutoMigrate(&models.User{}, &models.Stock{}, &models.StockPrice{}, &models.Reward{})
}
