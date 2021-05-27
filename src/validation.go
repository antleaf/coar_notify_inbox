package main

import (
	"github.com/piprate/json-gold/ld"
)

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
