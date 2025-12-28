package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/models"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func StartPriceWorker() {
	updateAllPrices()
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			logrus.Info("Starting stock prices update")
			updateAllPrices()
			logrus.Info("Successfully updated stock prices")
		}
	}()
}

func generateRandomPrice() decimal.Decimal {
	min := 250.0
	max := 500.0
	res := min + rand.Float64()*(max-min)
	return decimal.NewFromFloat(res).Round(4)
}

func updateAllPrices() {
	var stocks []models.Stock
	if err := database.DB.Where("is_active = ?", true).Find(&stocks).Error; err != nil {
		fmt.Printf("Failed to fetch stocks for price updates: %v\n", err)
		return
	}
	for _, v := range stocks {
		newPrice := generateRandomPrice()
		priceEntry := models.StockPrice{
			StockSymbol: v.Symbol,
			PriceINR:    newPrice,
			CapturedAt:  time.Now(),
		}
		logrus.WithFields(logrus.Fields{"StockSymbol": priceEntry.StockSymbol, "PriceInr": newPrice, "CapturedAt": time.Now()}).Infof("Adding updated price for stock: %v\n", v.Symbol)
		price := database.DB.Create(&priceEntry)
		if price.Error != nil {
			fmt.Printf("Failed to insert price for %s: %v\n", v.Symbol, price.Error)
			continue
		}
		logrus.Infof("Successfully updated price of stock: %v\n", v.Symbol)
	}
	logrus.Info("Stock prices updated successfully\n")
}
