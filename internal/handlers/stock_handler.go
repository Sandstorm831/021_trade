package handlers

import (
	"time"

	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateStock(c *gin.Context) {
	logrus.Info("Entering CreateStock handler.")
	var stockName struct {
		Symbol string
		Name   string
		Price  decimal.Decimal
	}
	if err := c.ShouldBindJSON(&stockName); err != nil {
		logrus.WithError(err).Error("Failed to bind JSON for CreateStock.")
		c.JSON(400, gin.H{
			"message": "Some error occurred",
		})
		return
	}
	stock := models.Stock{Symbol: stockName.Symbol, Name: stockName.Name}
	stockPrice := models.StockPrice{StockSymbol: stockName.Symbol, PriceINR: stockName.Price, CapturedAt: time.Now()}
	logrus.WithFields(logrus.Fields{
		"symbol": stockName.Symbol,
		"name":   stockName.Name,
		"price":  stockName.Price,
	}).Info("Attempting to create new stock and record its price.")

	database.DB.Transaction(func(tx *gorm.DB) error {
		stockRes := tx.Create(&stock)
		if stockRes.Error != nil {
			logrus.WithField("stockSymbol", stock.Symbol).WithError(stockRes.Error).Error("Some error occurred while creating stock in database.")
			c.JSON(400, gin.H{
				"message": "some error while creating stock in database",
			})
			return stockRes.Error
		}
		priceRes := tx.Create(&stockPrice)
		if priceRes.Error != nil {
			logrus.WithField("stockSymbol", stockPrice.StockSymbol).WithError(priceRes.Error).Error("Some error occurred while inserting stock price in database.")
			c.JSON(400, gin.H{
				"message": "some error while inserting stock price in database",
			})
			return priceRes.Error
		}
		return nil
	})

	logrus.Infof("Successfully created stock %s with price %s.", stock.Symbol, stockPrice.PriceINR)
	c.JSON(200, gin.H{
		"Stock": stock,
		"Price": stockPrice.PriceINR,
	})
}
