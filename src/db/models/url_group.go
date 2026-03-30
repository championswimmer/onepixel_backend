package models

import "gorm.io/gorm"

type UrlGroup struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	ShortPath string `gorm:"column:name;not null;size:10;uniqueIndex:idx_url_groups_short_path"`
	CreatorID uint64 `gorm:"not null"`
	Creator   User   `gorm:"foreignKey:CreatorID"`
}

func (UrlGroup) TableName() string {
	return "url_groups"
}
