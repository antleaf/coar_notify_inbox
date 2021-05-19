package main

const ldpContextURI = "http://www.w3.org/ns/ldp"

type Inbox struct {
	Notifications []Notification
}

func (inbox *Inbox) Populate() {
	inbox.Notifications = nil //TODO: is this necessary, or does the next line invoke garbage collection anyway?
	inbox.Notifications = make([]Notification, 0)
	rows, _ := db.Model(&NotificationDbRecord{}).Rows()
	defer rows.Close()
	for rows.Next() {
		var notificationDbRecord NotificationDbRecord
		db.ScanRows(rows, &notificationDbRecord)
		inbox.Notifications = append(inbox.Notifications, LoadNotificationFromDbRecord(notificationDbRecord))
	}
}

func (inbox *Inbox) Add(notification Notification) {
	inbox.Notifications = append(inbox.Notifications, notification)
}

func (inbox *Inbox) GetAsMapToPassToJsonRender() map[string]interface{} {
	inboxPayload := make(map[string]interface{})
	notificationURIs := make([]string, 0)
	for _, notification := range inbox.Notifications {
		notificationURIs = append(notificationURIs, notification.Url())
	}
	inboxPayload["@context"] = ldpContextURI
	inboxPayload["@id"] = site.InboxUrl()
	inboxPayload["contains"] = notificationURIs
	return inboxPayload
}
