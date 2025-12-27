package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Reward struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID       `gorm:"not null;index"`
	User           User            `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	StockSymbol    string          `gorm:"not null;index"`
	Stock          Stock           `gorm:"foreignKey:StockSymbol;references:Symbol;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Quantity       decimal.Decimal `gorm:"type:numeric(18,6);not null"`
	IdempotencyKey string          `gorm:"uniqueIndex;not null"`
	RewardedAt     time.Time       `gorm:"not null"`
}
