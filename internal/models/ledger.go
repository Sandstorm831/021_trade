package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LedgerEntry struct {
	ID          uint64          `gorm:"primaryKey"`
	RewardID    uuid.UUID       `gorm:"not null;index"`
	Reward      Reward          `gorm:"foreignKey:RewardID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StockSymbol *string         `gorm:"index"`
	Stock       *Stock          `gorm:"foreignKey:StockSymbol;references:Symbol;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AccountType string          `gorm:"type:varchar(50);not null"`
	Debit       decimal.Decimal `gorm:"type:numeric(18,6);default:0"`
	Credit      decimal.Decimal `gorm:"type:numeric(18,6);default:0"`
	Currency    string          `gorm:"type:varchar(10);not null"`
	CreatedAt   time.Time       `gorm:"not null"`
}
