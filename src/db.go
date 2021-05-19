package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type NotificationDbRecord struct {
	gorm.Model
	Payload string
	Sender  string
}

func InitialiseDb(dbPath string) error {
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		zapLogger.Error(err.Error())
		return err
	}
	db.AutoMigrate(&NotificationDbRecord{})
	return err
}
