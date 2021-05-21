package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
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
	//cliTesting2()
}

func cliTesting2() {
	generateUUID()
}

func cliTesting() {
	var err error
	err = testxpandJsonFromFile("data/payload_valid_notify.json")
	//err = testExpandJsonFromDbId(2)
	if err != nil {
		zapLogger.Error(err.Error())
	}
}

func runServer() {
	http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
}
