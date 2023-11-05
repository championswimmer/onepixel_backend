package models

import "gorm.io/gorm"

// Url db entity
type Url struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ShortURL  string `gorm:"index:short_url,unique,not null"`
	LongURL   string `gorm:"not null"`
	GroupID   uint
	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID"`
}

func (Url) TableName() string {
	return "urls"
}
