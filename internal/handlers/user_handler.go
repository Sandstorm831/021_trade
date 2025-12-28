package handlers

import (
	"math/rand"

	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RandomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	res := make([]byte, length)
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}
	return string(res)
}

func CreateUser(c *gin.Context) {
	var userName struct {
		Name string
	}
	if err := c.ShouldBindJSON(&userName); err != nil {
		c.JSON(400, gin.H{
			"message": "Some error occurred",
		})
		return
	}
	if userName.Name == "" {
		userName.Name = RandomString(8)
	}
	user := models.User{Name: userName.Name}
	logrus.Infof("Creating User: %s", user.Name)
	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "some error while inserting user in database",
		})
	}
	logrus.Infof("User %s, created successfuly", user.Name)
	logrus.WithFields(logrus.Fields{"Name": user.Name, "CreatedAt": user.CreatedAt})
	c.JSON(200, gin.H{
		"User": user,
	})
}
