package models

type Stock struct {
	Symbol   string `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`
}
