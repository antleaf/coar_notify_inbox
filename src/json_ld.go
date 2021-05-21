package main

import (
	"github.com/piprate/json-gold/ld"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"time"
)

func testExpandJsonFromNotificationRecord(notification Notification) error {
	var err error
	jsonLdProcessor := ld.NewJsonLdProcessor()
	var jsonLdProcessorOptions = ld.NewJsonLdOptions("")
	_, err = jsonLdProcessor.Expand(notification.Payload, jsonLdProcessorOptions)
	if err != nil {
		zapLogger.Error(err.Error())
		return err
	}
	return err
}

func testExpandJsonFromDbId(id uuid.UUID) error {
	var err error
	notification := LoadNotificationFromDbById(id)
	err = testExpandJsonFromNotificationRecord(notification)
	return err
}

func testxpandJsonFromFile(filePath string) error {
	var err error
	json, err := ioutil.ReadFile(filePath)
	if err != nil {
		zapLogger.Error(err.Error())
		return err
	}
	notification, err := NewNotification("", time.Now(), json)
	if err != nil {
		zapLogger.Error(err.Error())
		return err
	}
	err = testExpandJsonFromNotificationRecord(*notification)
	return err
}
