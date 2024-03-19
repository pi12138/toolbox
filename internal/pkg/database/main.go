package db

import (
	"database/sql"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var defaultDB *gorm.DB

func init() {
	var err error
	defaultDB, err = gorm.Open(sqlite.Open("toolbox.db"), &gorm.Config{})
	if err != nil {
		log.Panicf(`gorm.Open(sqlite.Open("toolbox.db"), &gorm.Config{}) error. %s`, err)
		return
	}
}

func D() *gorm.DB {
	return defaultDB
}

func Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return defaultDB.Transaction(fc, opts...)
}
