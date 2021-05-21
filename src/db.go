package main

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

type NotificationDbRecord struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	Payload    string
	ActivityId string
	Sender     string
}

//func (notificationDbRecord *NotificationDbRecord) BeforeCreate(scope *gorm.DB) error {
//	uuid := uuid.NewV4()
//	return scope.Set("ID", uuid).Error
//}

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
