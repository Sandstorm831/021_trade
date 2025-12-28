package main

import (
	"github.com/Sandstorm831/021_trade/internal/config"
	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/handlers"
	"github.com/Sandstorm831/021_trade/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	database.MigrateToDB()
	services.StartPriceWorker()
}

func main() {
	logrus.Info("Starting Gin Router")
	router := gin.Default()
	router.POST("/create-user", handlers.CreateUser)
	router.POST("/create-stock", handlers.CreateStock)
	router.POST("/reward", handlers.RecordReward)
	router.GET("/today-stocks/:userId", handlers.GetTodayRewards)
	router.GET("/historical-inr/:userId", handlers.GetHistoricalINR)
	router.GET("/stats/:userId", handlers.GetUserStats)
	router.GET("/portfolio/:userId", handlers.GetPortfolio)
	router.Run() // listens on 0.0.0.0:8080 by default
}
