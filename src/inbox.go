package main

const ldpContextURI = "http://www.w3.org/ns/ldp"

type Inbox struct {
	Notifications []Notification
}

func (inbox *Inbox) Initialise() {
	inbox.Notifications = make([]Notification, 0)
}

func (inbox *Inbox) SaveToDb() error {
	var err error
	for _, notification := range inbox.Notifications {
		notification.SaveToDb()
	}
	return err
}

func (inbox *Inbox) LoadFromDb() error {
	var err error
	rows, err := db.Model(&Notification{}).Rows()
	defer rows.Close()
	if err != nil {
		if err != nil {
			zapLogger.Error(err.Error())
			return err
		}
	}
	for rows.Next() {
		var notification Notification
		err = db.ScanRows(rows, &notification)
		if err != nil {
			zapLogger.Error(err.Error())
			return err
		}
		inbox.Notifications = append(inbox.Notifications, notification)
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
