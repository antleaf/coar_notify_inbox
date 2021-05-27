package main

const ldpContextURI = "http://www.w3.org/ns/ldp"

type Inbox struct {
	Notifications []Notification
}

func NewInbox() Inbox {
	var i = Inbox{}
	i.Notifications = make([]Notification, 0)
	db.Find(&i.Notifications)
	return i
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
