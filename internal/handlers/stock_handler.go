package handlers

import (
	"time"

	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func CreateStock(c *gin.Context) {
	var stockName struct {
		Symbol string
		Name   string
		Price  decimal.Decimal
	}
	if err := c.ShouldBindJSON(&stockName); err != nil {
		c.JSON(400, gin.H{
			"message": "Some error occurred",
		})
		return
	}
	stock := models.Stock{Symbol: stockName.Symbol, Name: stockName.Name}
	stockPrice := models.StockPrice{StockSymbol: stockName.Symbol, PriceINR: stockName.Price, CapturedAt: time.Now()}
	database.DB.Transaction(func(tx *gorm.DB) error {
		stockRes := database.DB.Create(&stock)
		if stockRes.Error != nil {
			c.JSON(400, gin.H{
				"message": "some error while creating stock in database",
			})
			return stockRes.Error
		}
		priceRes := database.DB.Create(&stockPrice)
		if priceRes.Error != nil {
			c.JSON(400, gin.H{
				"message": "some error while inserting stock price in database",
			})
			return priceRes.Error
		}
		return nil
	})

	c.JSON(200, gin.H{
		"Stock": stock,
		"Price": stockPrice.PriceINR,
	})
}
