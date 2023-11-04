package models

import (
	"gorm.io/gorm"
)

// User db entity
type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
