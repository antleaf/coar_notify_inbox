package main

import (
	"bytes"
	"encoding/json"
)

func (notification *Notification) Validate() error {
	var err error
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, notification.Payload)
	return err
}
