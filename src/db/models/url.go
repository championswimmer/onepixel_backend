package models

import "gorm.io/gorm"

// Url db entity
type Url struct {
	gorm.Model
	ID         uint64   `gorm:"primaryKey;autoIncrement:false"`
	ShortURL   string   `gorm:"not null;size:10;uniqueIndex:idx_urls_group_shortcode,priority:2"`
	LongURL    string   `gorm:"not null"`
	CreatorID  uint64   `gorm:"not null"`
	Creator    User     `gorm:"foreignKey:CreatorID"`
	UrlGroupID uint64   `gorm:"primaryKey;not null;default:0;uniqueIndex:idx_urls_group_shortcode,priority:1"`
	UrlGroup   UrlGroup `gorm:"foreignKey:UrlGroupID"`
}

func (Url) TableName() string {
	return "urls"
}
