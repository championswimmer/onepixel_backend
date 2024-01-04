package models

import "gorm.io/gorm"

type UrlGroup struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	Name      string `gorm:"unique,not null,size:10"`
	CreatorID uint   `gorm:"not null"`
	Creator   User   `gorm:"foreignKey:CreatorID"`
}

func (UrlGroup) TableName() string {
	return "url_groups"
}
