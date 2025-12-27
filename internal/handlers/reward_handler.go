package handlers

import (
	"time"

	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func ParseTimestamp(s string) (time.Time, error) {
	if s == "" {
		return time.Now().UTC(), nil
	}
	layoutWithTZ := "2006-01-02 15:04:05.999999-07:00"
	layoutNoTZ := "2006-01-02 15:04:05.999999"
	defaultLocation := time.UTC

	if t, err := time.Parse(layoutWithTZ, s); err == nil {
		return t, err
	}

	return time.ParseInLocation(layoutNoTZ, s, defaultLocation)
}

func RecordReward(c *gin.Context) {
	var rewardRecord struct {
		UserID         string
		StockSymbol    string
		Quantity       decimal.Decimal
		IdempotencyKey string
		RewardedAt     string
	}
	if err := c.ShouldBindBodyWithJSON(&rewardRecord); err != nil {
		c.JSON(400, gin.H{
			"message": "Some error occurred",
		})
		return
	}
	var userID uuid.UUID
	var rewardTime time.Time
	if id, err := uuid.Parse(rewardRecord.UserID); err != nil {
		c.JSON(400, gin.H{
			"message": "User Id is invalid",
		})
		return
	} else {
		userID = id
	}
	if pt, err := ParseTimestamp(rewardRecord.RewardedAt); err != nil {
		c.JSON(400, gin.H{
			"message": "Time formatting is not right",
		})
		return
	} else {
		rewardTime = pt
	}
	if rewardRecord.RewardedAt == "" {
		timeLayout := "2006-01-02 15:04:05.999999+00:00"
		rewardRecord.RewardedAt = time.Now().UTC().Format(timeLayout)
	}
	reward := models.Reward{
		UserID:         userID,
		StockSymbol:    rewardRecord.StockSymbol,
		Quantity:       rewardRecord.Quantity,
		IdempotencyKey: rewardRecord.IdempotencyKey,
		RewardedAt:     rewardTime,
	}
	res := database.DB.Create(&reward)
	if res.Error != nil {
		c.JSON(400, gin.H{
			"message": "Some error occurred while recording reward",
		})
		return
	}

	c.JSON(200, gin.H{
		"Reward": reward,
	})
}
