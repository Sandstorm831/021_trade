package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/Sandstorm831/021_trade/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
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
	reward := models.Reward{
		UserID:         userID,
		StockSymbol:    rewardRecord.StockSymbol,
		Quantity:       rewardRecord.Quantity,
		IdempotencyKey: rewardRecord.IdempotencyKey,
		RewardedAt:     rewardTime,
	}
	database.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&reward).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "Some error occurred while recording reward",
			})
			return err
		}
		var latestPrice models.StockPrice
		priceRes := tx.Where("stock_symbol = ?", rewardRecord.StockSymbol).Order("captured_at desc").First(&latestPrice)
		if priceRes.Error != nil {
			c.JSON(400, gin.H{
				"message": "Price for the Stock is not found",
			})
			return priceRes.Error
		}
		stockCostINR := rewardRecord.Quantity.Mul(latestPrice.PriceINR)
		feeINR := decimal.NewFromFloat(4.0 + rand.Float64()*10)
		totalCompanyCost := stockCostINR.Add(feeINR)
		ledgerEntryStock := models.LedgerEntry{
			RewardID:    reward.ID,
			StockSymbol: &rewardRecord.StockSymbol,
			AccountType: "USER_STOCK",
			Debit:       rewardRecord.Quantity,
			Currency:    "UNIT",
			CreatedAt:   rewardTime,
		}
		ledgerEntryCash := models.LedgerEntry{
			RewardID:    reward.ID,
			AccountType: "COMPANY_CASH",
			Credit:      totalCompanyCost,
			Currency:    "INR",
			CreatedAt:   rewardTime,
		}
		ledgerEntryFee := models.LedgerEntry{
			RewardID:    reward.ID,
			AccountType: "FEES_EXPENSE",
			Debit:       feeINR,
			Currency:    "INR",
			CreatedAt:   rewardTime,
		}
		if err := tx.Create(&ledgerEntryStock).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "Some error occurred while recording stock ledger entry",
			})
			return err
		}
		if err := tx.Create(&ledgerEntryCash).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "Some error occurred while recording company cash ledger entry",
			})
			return err
		}
		if err := tx.Create(&ledgerEntryFee).Error; err != nil {
			c.JSON(400, gin.H{
				"message": "Some error occurred while recording fee expense ledger entry",
			})
			return err
		}
		c.JSON(200, gin.H{
			"Reward":             reward,
			"stockLedgerEntryID": ledgerEntryStock.ID,
			"cashLedgerEntryID":  ledgerEntryCash.ID,
			"feeLedgerEntryID":   ledgerEntryFee.ID,
		})
		return nil
	})

}

func GetTodayRewards(c *gin.Context) {
	userId := c.Param("userId")
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var todayRewards []models.Reward
	if err := database.DB.Where("user_id = ? AND rewarded_at >= ?", userId, todayStart).Order("rewarded_at DESC").Find(&todayRewards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occurred while fetching the rewards",
		})
		return
	}
	c.JSON(http.StatusOK, todayRewards)
}
