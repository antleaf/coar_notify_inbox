package main

import (
	"github.com/robfig/cron/v3"
)

var cronRunner = cron.New()

func initialisePeriodicDbPersistence(period string) {
	cronRunner.AddFunc(period, func() {
		inbox.SaveToDb()
		zapLogger.Info("Saved Inbox to Database")
	})
	cronRunner.Start()
}
