package main

import (
	"fmt"
	_ "github.com/alecthomas/chroma/formatters"
	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

func Options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, HEAD, OPTIONS, POST")
	w.Header().Set("Accept-Post", "application/ld+json")
	pageRender.Text(w, http.StatusOK, "")
}

func HomePageGet(w http.ResponseWriter, r *http.Request) {
	var page = NewPage()
	page.Title = "Home"
	pageBytes, _ := GetPageBodyAsByteSliceFromFs("home.md")
	html, _ := GetHTMLFromMarkdown(pageBytes)
	page.Body = html
	pageRender.HTML(w, http.StatusOK, "page", page)
}

func InboxPost(w http.ResponseWriter, r *http.Request) {
	//TODO: Refactor to allow notifications to be captured even if they are bogus - capture HTTP headers and response codes. Do validation separately
	var err error
	notification := NewNotification(GetIP(r), time.Now())
	requestHeaderAsBytes, _ := yaml.Marshal(r.Header)
	notification.HttpRequest = string(requestHeaderAsBytes)
	payloadJson, err := ioutil.ReadAll(r.Body)
	if handlePostErrorCondition(err, w, 400, "Unable to read posted content", notification) {
		return
	}
	notification.Payload = payloadJson
	err = notification.Validate()
	if handlePostErrorCondition(err, w, 400, "Unable to parse posted content (must be JSON-LD)", notification) {
		return
	}
	notification.GeneratePayloadStructFromBytes()
	var page = NewPage()
	page.Params["notificationUrl"] = fmt.Sprint(notification.ID)
	page.Title = "Notification Response"
	w.Header().Set("Location", notification.Url())
	err = pageRender.HTML(w, http.StatusCreated, "post_success", page)
	if handlePostErrorCondition(err, w, 500, "Unable to process request", notification) {
		return
	}
	notification.HttpResponseCode = 201
	responseHeaderAsBytes, _ := yaml.Marshal(w.Header())
	notification.HttpResponseHeader = string(responseHeaderAsBytes)
	inbox.Add(*notification)
}

func handlePostErrorCondition(err error, w http.ResponseWriter, code int, message string, notification *Notification) bool {
	if err != nil {
		zapLogger.Error(err.Error())
		http.Error(w, message, code)
		notification.HttpResponseCode = code
		responseHeaderAsBytes, _ := yaml.Marshal(w.Header())
		notification.HttpResponseHeader = string(responseHeaderAsBytes)
		inbox.Add(*notification)
		return true
	} else {
		return false
	}
}

func InboxGet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "application/json" {
		pageRender.JSON(w, http.StatusOK, inbox.GetAsMapToPassToJsonRender())
	} else {
		page := NewInboxPage(inbox.Notifications)
		page.Title = "Inbox"
		pageRender.HTML(w, http.StatusOK, "inbox", page)
	}
}

func InboxGetJSON(w http.ResponseWriter, r *http.Request) {
	pageRender.JSON(w, http.StatusOK, inbox.GetAsMapToPassToJsonRender())
}

func InboxNotificationGet(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, _ := uuid.FromString(idString)
	notification := LoadNotificationFromDbById(id)
	var page = NewNotificationPage(notification)
	page.Title = "Notification"
	pageRender.HTML(w, http.StatusOK, "notification", page)
}

func GetPageBodyAsByteSliceFromFs(filename string) ([]byte, error) {
	var err error
	pageBodyBytes, err := pageAssets.ReadFile(filepath.Join("pages", filename))
	return pageBodyBytes, err
}

func GetIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
