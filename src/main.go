package main

import (
	"flag"
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
	debugPtr := flag.Bool("debug", false, "Enable debug logging")
	portPtr := flag.Int("port", 1313, "Port number")
	dbPathPtr := flag.String("db", "./db.sqlite", "Path to to Database file")
	baseUrlPtr := flag.String("baseUrl", "http://localhost:1313", "Base URL")
	flag.Parse()
	config.initialise(debugPtr, portPtr, dbPathPtr, baseUrlPtr)
	err := InitialiseDb(config.DbFilePath)
	inbox.Populate()
	initialiseRendering()
	if err != nil {
		zapLogger.Fatal(err.Error())
	}
	router = ConfigureRouter()
	http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
}
