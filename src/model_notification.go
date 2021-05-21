package main

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"time"
)

type Notification struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	Sender        string
	Payload       []byte
	PayloadStruct map[string]interface{} `gorm:"-"`
	ActivityId    string
}

func NewNotification(sender string, timestamp time.Time, payloadJson []byte) (*Notification, error) {
	var n Notification
	var err error
	n.ID = uuid.NewV4()
	n.Sender = sender
	n.CreatedAt = timestamp
	n.Payload = payloadJson
	err = n.GeneratePayloadStructFromBytes()
	if err != nil {
		zapLogger.Error(err.Error())
		return &n, err
	}
	n.ActivityId = n.ExtractActivityId()
	return &n, err
}

func (notification *Notification) GeneratePayloadStructFromBytes() error {
	var err error
	err = json.Unmarshal(notification.Payload, &(notification.PayloadStruct))
	if err != nil {
		zapLogger.Error(err.Error())
		return err
	}
	return err
}

func (notification *Notification) SaveToDb() error {
	var err error
	tempNotification := Notification{} // use this so the db lookup does not overwrite the incoming notification with the saved one
	err = db.First(&tempNotification, notification.ID).Error
	if err == nil {
		err = db.Save(&notification).Error
	} else {
		err = db.Create(&notification).Error
	}
	if err != nil {
		zapLogger.Error(err.Error())
	}
	return err
}

func LoadNotificationFromDbById(id uuid.UUID) Notification {
	notification := Notification{}
	db.First(&notification, id)
	notification.GeneratePayloadStructFromBytes()
	return notification
}

func (notification *Notification) Url() string {
	return site.InboxUrl() + notification.ID.String()
}

func (notification *Notification) ExtractActivityId() string {
	return fmt.Sprintf("%v", notification.PayloadStruct["id"])
}

func (notification *Notification) FormattedTimestamp() string {
	return notification.CreatedAt.Format("2006-01-02 15:04:05")
}

func (notification *Notification) HTMLFormattedPayload() template.HTML {
	payloadBytes, err := json.MarshalIndent(notification.PayloadStruct, "", "    ")
	if err != nil {
		zapLogger.Error(err.Error())
	}
	payloadJson := fmt.Sprintf("```json\n%s\n", payloadBytes)
	htmlPayload, _ := GetHTMLFromMarkdown([]byte(payloadJson))
	return htmlPayload
}

//type PayloadStruct struct {
//	Id     string   `json:"id"`
//	Type   []string `json:"type"`
//	Actor  Actor    `json:"actor, omitempty"`
//	Origin Service  `json:"origin"`
//	Target Service  `json:"origin"`
//	Object *Object  `json:"object,omitempty"`
//}
//
//func (payload *PayloadStruct) UnMarshall(payloadJson []byte) {
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
