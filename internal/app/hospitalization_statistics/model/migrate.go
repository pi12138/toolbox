package model

import "gorm.io/gorm"

func Migrate(tx *gorm.DB) {
	tx.AutoMigrate(&Item{})
	tx.AutoMigrate(&Daily{})
}
