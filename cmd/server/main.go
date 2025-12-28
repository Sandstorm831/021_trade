package main

import (
	"fmt"

	"github.com/Sandstorm831/021_trade/internal/config"
	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/handlers"
	"github.com/Sandstorm831/021_trade/internal/services"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	database.MigrateToDB()
	services.StartPriceWorker()
}

func main() {

	fmt.Println("Hello World")
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/create-user", handlers.CreateUser)
	router.POST("/create-stock", handlers.CreateStock)
	router.POST("/reward", handlers.RecordReward)
	router.GET("/today-stocks/:userId", handlers.GetTodayRewards)
	router.GET("/historical-inr/:userId", handlers.GetHistoricalINR)
	router.Run() // listens on 0.0.0.0:8080 by default
}
