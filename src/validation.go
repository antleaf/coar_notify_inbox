package main

import (
	"github.com/piprate/json-gold/ld"
)

//func (notification *Notification) CheckPayloadIsWellFormedJson() error {
//	var err error
//	buffer := new(bytes.Buffer)
//	err = json.Compact(buffer, notification.Payload)
//	if err != nil {
//		zapLogger.Error(err.Error())
//		return err
//	}
//	return err
//}

func (notification *Notification) CheckPayloadIsJsonLd() error {
	var err error
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	payloadAsInterface, err := notification.ExpressPayloadAsInterface()
	if err != nil {
		zapLogger.Error(err.Error())
		notification.AddToProcessLog(err.Error())
		return err
	}
	_, err = proc.Expand(payloadAsInterface, options)
	if err != nil {
		zapLogger.Error(err.Error())
		notification.AddToProcessLog(err.Error())
		return err
	}
	notification.AddToProcessLog("Appears to be valid JSON-LD")
	return err
}

//func testExpandJsonFromNotificationRecord(notification Notification) error {
//	var err error
//	jsonLdProcessor := ld.NewJsonLdProcessor()
//	var jsonLdProcessorOptions = ld.NewJsonLdOptions("")
//	_, err = jsonLdProcessor.Expand(notification.PayloadStruct, jsonLdProcessorOptions)
//	if err != nil {
//		zapLogger.Error(err.Error())
//		return err
//	}
//	return err
//}
//
//func testExpandJsonFromDbId(id uuid.UUID) error {
//	var err error
//	notification := LoadNotificationFromDbById(id)
//	err = testExpandJsonFromNotificationRecord(notification)
//	return err
//}

//func testxpandJsonFromFile(filePath string) error {
//	var err error
//	err = ioutil.ReadFile(filePath)
//	if err != nil {
//		zapLogger.Error(err.Error())
//		return err
//	}
//	notification := NewNotification("", time.Now())
//	err = testExpandJsonFromNotificationRecord(*notification)
//	return err
//}
