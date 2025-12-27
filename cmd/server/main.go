package main

import (
	"fmt"

	"github.com/Sandstorm831/021_trade/internal/config"
	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/handlers"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnvVariables()
	database.ConnectToDB()
	database.MigrateToDB()
}

func main() {

	fmt.Println("Hello World")
	router := gin.Default()
	router.GET("/", handlers.CreateUser)
	router.Run() // listens on 0.0.0.0:8080 by default
}
