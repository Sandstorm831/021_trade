package handlers

import (
	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateStock(c *gin.Context) {
	var stockName struct {
		Symbol string
		Name   string
	}
	if err := c.ShouldBindJSON(&stockName); err != nil {
		c.JSON(400, gin.H{
			"message": "Some error occurred",
		})
		return
	}
	stock := models.Stock{Symbol: stockName.Symbol, Name: stockName.Name}
	result := database.DB.Create(&stock)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "some error while creating stock in database",
		})
	}

	c.JSON(200, gin.H{
		"Stock": stock,
	})
}
