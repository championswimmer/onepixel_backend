package models

import (
	"gorm.io/gorm"
)

// User db entity
type User struct {
	gorm.Model
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Verified bool   `gorm:"default:false"`
}

func (User) TableName() string {
	return "users"
}
