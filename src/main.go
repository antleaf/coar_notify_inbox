package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"
	"go.uber.org/zap"
	"net/http"
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

	runServer()
	//cliTesting()
}

func cliTesting() {
	id, _ := uuid.FromString("d799534f-2c5e-40ae-a3db-325703c1e3ec")
	notification := LoadNotificationFromDbById(id)
	zapLogger.Debug(notification.ActivityId)
}

func runServer() {
	http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
}
