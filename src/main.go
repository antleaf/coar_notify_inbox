package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

var zapLogger *zap.Logger
var router chi.Router
var config = Config{}
var pageRender *render.Render
var inbox = Inbox{}

func main() {
	var err error
	err = configure()
	if err == nil {
		zapLogger.Info("System configured OK")
	} else {
		zapLogger.Fatal("Aborting start-up due to configuration error")
	}

	//http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
	err = testxpandJsonFromFile("data/payload_json_not_valid.json")
	//err = testxpandJsonFromFile("data/payload_valid_notify.json")
	//err = testxpandJsonFromFile("data/payload_not_json.txt")
	//notification := LoadNotificationFromDbById(2)
	//err = testxpandJsonFromInterfaces(notification.Payload)
	if err != nil {
		zapLogger.Error(err.Error())
	}
}
