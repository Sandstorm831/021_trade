package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Stock struct {
	Symbol   string `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`
}

type StockPrice struct {
	ID          uint64          `gorm:"primaryKey"`
	StockSymbol string          `gorm:"not null;index:idx_stock_time"`
	Stock       Stock           `gorm:"foreignKey:StockSymbol;references:Symbol;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PriceINR    decimal.Decimal `gorm:"type:numeric(18,4);not null"`
	CapturedAt  time.Time       `gorm:"index:idx_stock_time,sort:desc;not null"`
}
