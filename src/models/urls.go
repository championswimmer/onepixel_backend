package models

import "gorm.io/gorm"

// Url db entity
type Url struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	ShortURL  string `gorm:"uniqueIndex,not null"`
	LongURL   string `gorm:"not null"`
	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID"`
}

func (Url) TableName() string {
	return "urls"
}
