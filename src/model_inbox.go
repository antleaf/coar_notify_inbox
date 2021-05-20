package main

const ldpContextURI = "http://www.w3.org/ns/ldp"

type Inbox struct {
	Notifications []Notification
}

func (inbox *Inbox) Populate() error {
	var err error
	inbox.Notifications = nil //TODO: is this necessary, or does the next line invoke garbage collection anyway?
	inbox.Notifications = make([]Notification, 0)
	rows, err := db.Model(&NotificationDbRecord{}).Rows()
	defer rows.Close()
	if err != nil {
		if err != nil {
			zapLogger.Error(err.Error())
			return err
		}
	}
	for rows.Next() {
		var notificationDbRecord NotificationDbRecord
		err = db.ScanRows(rows, &notificationDbRecord)
		if err != nil {
			zapLogger.Error(err.Error())
			return err
		}
		inbox.Notifications = append(inbox.Notifications, LoadNotificationFromDbRecord(notificationDbRecord))
	}
	return err
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
