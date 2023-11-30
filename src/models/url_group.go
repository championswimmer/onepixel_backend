package models

type UrlGroup struct {
	ID        string `gorm:"primaryKey;autoIncrement:false"`
	Name      string `gorm:"unique,not null,size:10"`
	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID"`
}
