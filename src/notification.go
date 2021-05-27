package main

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"time"
)

type Notification struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time `sql:"index"`
	Sender             string
	Payload            []byte
	ActivityId         string
	HttpRequest        string
	HttpResponseHeader string
	HttpResponseCode   int
	ProcessLog         string
}

func NewNotification(sender string, timestamp time.Time) *Notification {
	var n Notification
	n.ID = uuid.NewV4()
	n.Sender = sender
	n.CreatedAt = timestamp
	return &n
}

func (notification *Notification) AddToProcessLog(message string) {
	notification.ProcessLog += message + "\n"
}

func (notification *Notification) ExpressPayloadAsMap() (map[string]interface{}, error) {
	var err error
	var payloadMap map[string]interface{}
	err = json.Unmarshal(notification.Payload, &(payloadMap))
	if err != nil {
		zapLogger.Error(err.Error())
	}
	return payloadMap, err
}

func (notification *Notification) ExpressPayloadAsInterface() (interface{}, error) {
	return notification.ExpressPayloadAsMap()
}

func (notification *Notification) SaveToDb() error {
	var err error
	err = db.First(&Notification{}, notification.ID).Error
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
	notification.ExpressPayloadAsMap()
	return notification
}

func (notification *Notification) Url() string {
	return site.InboxUrl() + notification.ID.String()
}

func (notification *Notification) ProcessPayload() {
	var err error
	payloadStruct, err := notification.ExpressPayloadAsMap()
	if err == nil {
		notification.ActivityId = fmt.Sprintf("%v", payloadStruct["id"])
	}
}

func (notification *Notification) FormattedTimestamp() string {
	return notification.CreatedAt.Format("2006-01-02 15:04:05")
}

func (notification *Notification) HTMLFormattedPayload() template.HTML {
	_, err := notification.ExpressPayloadAsMap()
	if err == nil {
		payloadJson := fmt.Sprintf("```json\n%s\n```\n", notification.Payload)
		htmlPayload, _ := GetHTMLFromMarkdown([]byte(payloadJson))
		return htmlPayload
	} else {
		htmlPayload, _ := GetHTMLFromMarkdown([]byte(notification.Payload))
		return htmlPayload
	}
}

func (notification *Notification) HTMLFormattedHttpHeaders() template.HTML {
	markdown := "#### HTTP Request Headers\n"
	markdown += fmt.Sprintf("```yaml\n%s\n```\n", notification.HttpRequest)
	markdown += "#### HTTP Response Headers\n"
	markdown += fmt.Sprintf("```yaml\n%s\n```\n", notification.HttpResponseHeader)
	htmlPayload, _ := GetHTMLFromMarkdown([]byte(markdown))
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
