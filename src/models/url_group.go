package models

type UrlGroup struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	Name      string `gorm:"unique,not null,size:10"`
	CreatorID uint   `gorm:"not null"`
	Creator   User   `gorm:"foreignKey:CreatorID"`
}
