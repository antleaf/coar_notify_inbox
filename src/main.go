package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/piprate/json-gold/ld"
	uuid "github.com/satori/go.uuid"
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

	runServer()
	//cliTesting()
}

func cliTesting() {
	id, _ := uuid.FromString("d799534f-2c5e-40ae-a3db-325703c1e3ec")
	notification := Notification{}
	db.First(&notification, id)
	zapLogger.Debug(notification.ActivityId)

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/n-quads"
	doc := map[string]interface{}{
		"@context": "https://json-ld.org/contexts/person.jsonld",
		"@id":      "http://dbpedia.org/resource/John_Lennon",
		"name":     "John Lennon",
		"born":     "1940-10-09",
		"spouse":   "http://dbpedia.org/resource/Cynthia_Lennon",
	}

	doc, _ = notification.ExpressPayloadAsMap()
	quads, err := proc.ToRDF(doc, options)

	if err != nil {
		zapLogger.Error(err.Error())
	} else {
		quadStrings := quads.(string)
		zapLogger.Debug(quadStrings)
	}

}

func runServer() {
	http.ListenAndServe(fmt.Sprintf(":%v", config.Port), router)
}
