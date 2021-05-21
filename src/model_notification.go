package main

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"time"
)

type Notification struct {
	ID        uuid.UUID
	Sender    string
	CreatedAt time.Time
	Payload   map[string]interface{}
}

func NewNotification(sender string, timestamp time.Time, payloadJson []byte) (*Notification, error) {
	var n Notification
	var err error
	n.ID = uuid.NewV4()
	n.Sender = sender
	n.CreatedAt = timestamp
	err = json.Unmarshal(payloadJson, &(n.Payload))
	if err != nil {
		zapLogger.Error(err.Error())
	}
	return &n, err
}

func (notification *Notification) SaveToDb() error {
	var err error
	notificationRecord := NotificationDbRecord{}
	notificationRecord.ID = notification.ID
	notificationRecord.Sender = notification.Sender
	notificationRecord.CreatedAt = notification.CreatedAt
	notificationRecord.ActivityId = notification.ActivityId()
	payloadBytes, err := json.Marshal(notification.Payload)
	if err != nil {
		zapLogger.Error(err.Error())
		return err
	}
	notificationRecord.Payload = string(payloadBytes)
	tempNotificationRecord := NotificationDbRecord{} // use this so the db lookup does not overwrite the incoming notification with the saved one
	err = db.First(&tempNotificationRecord, notification.ID).Error
	if err == nil {
		err = db.Save(&notificationRecord).Error
	} else {
		err = db.Create(&notificationRecord).Error
		notification.ID = notificationRecord.ID
	}
	if err != nil {
		zapLogger.Error(err.Error())
	}
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

func LoadNotificationFromDbById(id uuid.UUID) Notification {
	notificationRecord := NotificationDbRecord{}
	db.First(&notificationRecord, id)
	notification := LoadNotificationFromDbRecord(notificationRecord)
	return notification
}

func (notification *Notification) Url() string {
	return site.InboxUrl() + notification.ID.String()
}

func (notification *Notification) ActivityId() string {
	return fmt.Sprintf("%v", notification.Payload["id"])
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
