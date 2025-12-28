package handlers

import (
	"fmt"
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

func GetHistoricalINR(c *gin.Context) {
	userID := c.Param("userId")
	const layout = "2006-01-02"
	// Finding rewards
	var rewards []models.Reward
	if err := database.DB.Where("user_id = ?", userID).Order("rewarded_at asc").Find(&rewards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occurred while fetching the rewards",
		})
		return
	}

	if len(rewards) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Finding prices
	type DailyPrice struct {
		Date        time.Time
		Price       decimal.Decimal
		StockSymbol string
	}
	var dailyClosingPrice []DailyPrice
	if err := database.DB.Raw(`
	SELECT DISTINCT ON (DATE(captured_at), stock_symbol)
		DATE(captured_at) as date, price_inr as price, stock_symbol
	FROM stock_prices
	WHERE captured_at < CURRENT_DATE
	ORDER BY DATE(captured_at), stock_symbol, captured_at DESC
	`).Find(&dailyClosingPrice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occurred while getting prices",
		})
		return
	}

	fmt.Println(dailyClosingPrice)

	firstDate := rewards[0].RewardedAt
	var closingPriceBeforeFirstDate []DailyPrice
	if err := database.DB.Raw(`
	SELECT DISTINCT ON (stock_symbol)
		DATE(captured_at) as date, price_inr as price, stock_symbol
	FROM stock_prices
	WHERE captured_at < ?
	ORDER BY stock_symbol, captured_at DESC
	`, firstDate.Format(layout)).Find(&closingPriceBeforeFirstDate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occurred while getting prices before first date",
		})
		return
	}
	fmt.Printf("Before first date: %v\n", closingPriceBeforeFirstDate)
	type HistoricalINR struct {
		Date  string
		Value decimal.Decimal
	}

	holdings := make(map[string]map[string]decimal.Decimal)
	priceReg := make(map[string]map[string]decimal.Decimal)
	pastDayStr := firstDate.AddDate(0, 0, -1).Format(layout)
	priceReg[pastDayStr] = make(map[string]decimal.Decimal)
	for _, v := range closingPriceBeforeFirstDate {
		priceReg[pastDayStr][v.StockSymbol] = v.Price
	}
	yesterday := time.Now().AddDate(0, 0, -1)
	rewardIndex := 0
	dailyClosingPriceIndex := 0
	var historicalVal []HistoricalINR

	// this for loop is for aggregating all the quantity and prices of stock present over the past
	for d := firstDate; !d.After(yesterday); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format(layout)

		// adding entries of StockMap for each day, where each StockMap is a map of stock_symbol and it's qty rewarded that day
		for rewardIndex < len(rewards) && rewards[rewardIndex].RewardedAt.Format(layout) == dateStr {
			if _, ok := holdings[dateStr]; !ok {
				holdings[dateStr] = make(map[string]decimal.Decimal)
			}
			if _, ok := holdings[dateStr][rewards[rewardIndex].StockSymbol]; !ok {
				holdings[dateStr][rewards[rewardIndex].StockSymbol] = rewards[rewardIndex].Quantity
			} else {
				holdings[dateStr][rewards[rewardIndex].StockSymbol] = holdings[dateStr][rewards[rewardIndex].StockSymbol].Add(rewards[rewardIndex].Quantity)
			}
			rewardIndex += 1
		}

		// cumulating the previous day stocks quantities if not the first day
		if dateStr != firstDate.Format(layout) {
			prevDay := d.AddDate(0, 0, -1).Format(layout)
			prevHoldingEntry := holdings[prevDay]
			for s_sym, s_qty := range prevHoldingEntry {
				if _, ok := holdings[dateStr]; !ok {
					holdings[dateStr] = make(map[string]decimal.Decimal)
				}
				if _, ok := holdings[dateStr][s_sym]; !ok {
					holdings[dateStr][s_sym] = s_qty
				} else {
					holdings[dateStr][s_sym] = holdings[dateStr][s_sym].Add(s_qty)
				}
			}
		}

		// adding entries of PriceMap for each day, where each PriceMap is a map of stock_symbol and it's price that day
		for dailyClosingPriceIndex < len(dailyClosingPrice) && dailyClosingPrice[dailyClosingPriceIndex].Date.Format(layout) == dateStr {
			if _, ok := priceReg[dateStr]; !ok {
				priceReg[dateStr] = make(map[string]decimal.Decimal)
			}
			priceReg[dateStr][dailyClosingPrice[dailyClosingPriceIndex].StockSymbol] = dailyClosingPrice[dailyClosingPriceIndex].Price
			dailyClosingPriceIndex += 1
		}

		// adding previous day stocks prices
		prevDay := d.AddDate(0, 0, -1).Format(layout)
		prevPriceEntry := priceReg[prevDay]
		for s_sym, s_price := range prevPriceEntry {
			if _, ok := priceReg[dateStr]; !ok {
				priceReg[dateStr] = make(map[string]decimal.Decimal)
			}
			if _, ok := priceReg[dateStr][s_sym]; !ok {
				priceReg[dateStr][s_sym] = s_price
			}
		}

		currDayVal := decimal.Zero
		for s_sym, s_qty := range holdings[dateStr] {
			s_price := priceReg[dateStr][s_sym]
			currDayVal = currDayVal.Add(s_qty.Mul(s_price))
			fmt.Printf("Date: %v, s_sym: %v, s_price: %v, s_qty: %v, CurrVal: %v\n", dateStr, s_sym, s_price, s_qty, currDayVal)
		}
		historicalVal = append(historicalVal, HistoricalINR{Date: dateStr, Value: currDayVal})
	}
	c.JSON(http.StatusOK, historicalVal)

}
