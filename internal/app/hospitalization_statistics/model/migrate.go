package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

func Migrate(tx *gorm.DB) {
	var tables []Tabler = []Tabler{
		&Item{},
		&Daily{},
	}

	for i := 0; i < len(tables); i++ {
		if err := tx.AutoMigrate(tables[i]); err != nil {
			fmt.Printf("table %s migrate error. %s\n", tables[i].TableName(), err)
		} else {
			fmt.Printf("table %s migrate success.\n", tables[i].TableName())
		}
	}
}
