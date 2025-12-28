package handlers

import (
	"net/http"
	"time"

	"github.com/Sandstorm831/021_trade/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func GetUserStats(c *gin.Context) {
	logrus.Info("Entering GetUserStats handler.")
	userID := c.Param("userId")

	type StockShares struct {
		Symbol   string          `json:"symbol"`
		Quantity decimal.Decimal `json:"total_quantity"`
	}

	type StatsResponse struct {
		TodayRewards []StockShares   `json:"today_rewards"`
		TotalValue   decimal.Decimal `json:"total_portfolio_value_inr"`
	}

	// Total shares rewarded
	var todayRewards []StockShares
	if err := database.DB.Raw(`
		SELECT stock_symbol as symbol, SUM(quantity) as quantity
		FROM rewards
		WHERE user_id = ? AND DATE(rewarded_at) >= DATE(?)
		GROUP BY stock_symbol
	`, userID, time.Now()).Scan(&todayRewards).Error; err != nil {
		logrus.WithField("userID", userID).WithError(err).Error("Some error occurred while fetching today's stats.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occurred while fetching today's stats",
		})
		return
	}
	logrus.Debugf("Today's rewards for user %s: %v", userID, todayRewards)

	type Holding struct {
		Symbol string
		Total  decimal.Decimal
	}
	var allHoldings []Holding
	if err := database.DB.Raw(`
    	SELECT stock_symbol AS symbol, SUM(quantity) AS total
    	FROM rewards
    	WHERE user_id = ?
    	GROUP BY stock_symbol
	`, userID).Scan(&allHoldings).Error; err != nil {
		logrus.WithField("userID", userID).WithError(err).Error("Some error occurred while fetching stock holdings for user stats.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occurred while fetching stock holdings",
		})
		return
	}
	logrus.Debugf("All holdings for user %s: %v", userID, allHoldings)

	type LatestPrice struct {
		StockSymbol string
		PriceINR    decimal.Decimal
	}
	var prices []LatestPrice
	database.DB.Raw(`
		SELECT DISTINCT ON (stock_symbol) stock_symbol, price_inr
		FROM stock_prices
		ORDER BY stock_symbol, captured_at DESC
	`).Scan(&prices)

	// Create a quick lookup map for prices
	priceMap := make(map[string]decimal.Decimal)
	for _, p := range prices {
		priceMap[p.StockSymbol] = p.PriceINR
	}

	// 3. Multiply and Sum
	totalPortfolioValue := decimal.Zero
	for _, h := range allHoldings {
		if price, ok := priceMap[h.Symbol]; ok {
			totalPortfolioValue = totalPortfolioValue.Add(h.Total.Mul(price))
		}
	}

	logrus.Infof("Successfully fetched user stats for user %s. Total portfolio value: %s", userID, totalPortfolioValue)
	// Final Response
	c.JSON(http.StatusOK, StatsResponse{
		TodayRewards: todayRewards,
		TotalValue:   totalPortfolioValue,
	})
}

func GetPortfolio(c *gin.Context) {
	logrus.Info("Entering GetPortfolio handler.")
	userID := c.Param("userId")

	type PortfolioItem struct {
		Symbol        string          `json:"symbol"`
		TotalQuantity decimal.Decimal `json:"total_quantity"`
		CurrentPrice  decimal.Decimal `json:"current_price"`
		CurrentValue  decimal.Decimal `json:"current_value"`
	}
	var portfolio []PortfolioItem

	type Holding struct {
		Symbol string
		Total  decimal.Decimal
	}
	var allHoldings []Holding
	if err := database.DB.Raw(`
    	SELECT stock_symbol AS symbol, SUM(quantity) AS total
    	FROM rewards
    	WHERE user_id = ?
    	GROUP BY stock_symbol
	`, userID).Scan(&allHoldings).Error; err != nil {
		logrus.WithField("userID", userID).WithError(err).Error("Some error occurred while fetching stock holdings for portfolio.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occurred while fetching stock holdings",
		})
		return
	}

	logrus.Debugf("All holdings for user %s: %v", userID, allHoldings)

	type LatestPrice struct {
		StockSymbol string
		PriceINR    decimal.Decimal
	}
	var prices []LatestPrice
	database.DB.Raw(`
		SELECT DISTINCT ON (stock_symbol) stock_symbol, price_inr
		FROM stock_prices
		ORDER BY stock_symbol, captured_at DESC
	`).Scan(&prices)

	// Create a quick lookup map for prices
	priceMap := make(map[string]decimal.Decimal)
	for _, p := range prices {
		priceMap[p.StockSymbol] = p.PriceINR
	}

	// 3. Multiply and Sum
	totalPortfolioValue := decimal.Zero
	for _, h := range allHoldings {
		if price, ok := priceMap[h.Symbol]; ok {
			item := PortfolioItem{
				Symbol:        h.Symbol,
				TotalQuantity: h.Total,
				CurrentPrice:  price,
				CurrentValue:  h.Total.Mul(price),
			}
			portfolio = append(portfolio, item)
			totalPortfolioValue = totalPortfolioValue.Add(h.Total.Mul(price))
		}
	}

	logrus.Infof("Successfully fetched portfolio for user %s. Total portfolio value: %s", userID, totalPortfolioValue)
	// Final Response
	c.JSON(http.StatusOK, gin.H{
		"Portfoli":   portfolio,
		"TotalValue": totalPortfolioValue,
	})
}
