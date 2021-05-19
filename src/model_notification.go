package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"
)

type Notification struct {
	ID        uint
	Sender    string
	CreatedAt time.Time
	Payload   map[string]interface{}
}

func NewNotification(sender string, timestamp time.Time, payloadJson []byte) (*Notification, error) {
	var n Notification
	var err error
	n.Sender = sender
	n.CreatedAt = timestamp
	err = json.Unmarshal(payloadJson, &(n.Payload))
	if err != nil {
		zapLogger.Error(err.Error())
	}
	return &n, err
}

func (notification *Notification) SaveNotificationToDb() error {
	var err error
	notificationRecord := NotificationDbRecord{}
	notificationRecord.Sender = notification.Sender
	notificationRecord.CreatedAt = notification.CreatedAt
	payloadBytes, err := json.Marshal(notification.Payload)
	if err != nil {
		zapLogger.Error(err.Error())
		return err
	}
	notificationRecord.Payload = string(payloadBytes)
	err = db.Create(&notificationRecord).Error
	if err != nil {
		zapLogger.Error(err.Error())
		return err
	}
	notification.ID = notificationRecord.ID
	//notification.Payload["id"] = notification.Url()
	//payloadBytes, err = json.Marshal(notification.Payload)
	//notificationRecord.Payload = string(payloadBytes)
	//db.Save(&notificationRecord)
	return err
}

func LoadNotificationFromDbRecord(notificationRecord NotificationDbRecord) Notification {
	notification := Notification{}
	notification.ID = notificationRecord.ID
	notification.Sender = notificationRecord.Sender
	notification.CreatedAt = notificationRecord.CreatedAt
	err := json.Unmarshal([]byte(notificationRecord.Payload), &(notification.Payload))
	if err != nil {
		zapLogger.Error(err.Error())
	}
	return notification
}

func LoadNotificationFromDbById(id uint) Notification {
	notificationRecord := NotificationDbRecord{}
	db.First(&notificationRecord, id)
	notification := LoadNotificationFromDbRecord(notificationRecord)
	return notification
}

func (notification *Notification) Url() string {
	return site.InboxUrl() + strconv.FormatUint(uint64(notification.ID), 10)
}

func (notification *Notification) FormattedTimestamp() string {
	return notification.CreatedAt.Format("2006-01-02 15:04:05")
}

func (notification *Notification) HTMLFormattedPayload() template.HTML {
	//payloadBytes,err := json.Marshal(notification.Payload)
	payloadBytes, err := json.MarshalIndent(notification.Payload, "", "    ")
	if err != nil {
		zapLogger.Error(err.Error())
	}
	payloadJson := fmt.Sprintf("```json\n%s\n", payloadBytes)
	htmlPayload, _ := GetHTMLFromMarkdown([]byte(payloadJson))
	return htmlPayload
}

//type Payload struct {
//	Id     string   `json:"id"`
//	Type   []string `json:"type"`
//	Actor  Actor    `json:"actor, omitempty"`
//	Origin Service  `json:"origin"`
//	Target Service  `json:"origin"`
//	Object *Object  `json:"object,omitempty"`
//}
//
//func (payload *Payload) UnMarshall(payloadJson []byte) {
//	err := json.Unmarshal([]byte(payloadJson), &payload)
//	zapLogger.Error(err.Error())
//}
//
//type Object struct {
//	Id    string `json:"id"`
//	Type  []string `json:"type"`
//	CiteAs string `json:"ietf:cite-as"`
//	Url *Url `json:"url,omitempty"`
//}
//
//type Url struct {
//	Id    string `json:"id"`
//	Type  []string `json:"type"`
//	MediaType string `json:"media-type"`
//}
//
//type Service struct {
//	Id    string `json:"id"`
//	Type  []string `json:"type"`
//	Inbox string `json:"ldp:inbox"`
//}
//
//type Actor struct {
//	Id       string `json:"id"`
//	Type     []string `json:"type"`
//	Name     string `json:"name"`
//	LdpInbox string `json:"ldp:inbox,omitempty"`
//}
//
