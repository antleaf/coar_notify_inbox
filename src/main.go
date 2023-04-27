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

func main() {
	var err error
	err = configure()
	if err == nil {
		zapLogger.Info("System configured OK")
	} else {
		zapLogger.Fatal("Aborting start-up due to configuration error")
	}
	initialisePasswordlessAuthentication()
	runServer()
}

func runServer() {
	http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
}
