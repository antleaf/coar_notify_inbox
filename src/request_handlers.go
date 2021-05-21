package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/alecthomas/chroma/formatters"
	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"
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

//func InboxPost(w http.ResponseWriter, r *http.Request) {
//	var err error
//	payloadJson, err := ioutil.ReadAll(r.Body)
//	if handleRequestProcessingError(w, err, "Unable to read posted content", 400) {
//		return
//	}
//	buffer := new(bytes.Buffer)
//	err = json.Compact(buffer, payloadJson)
//	if handleRequestProcessingError(w, err, "Unable to parse posted content (must be JSON-LD)", 400) {
//		return
//	}
//	payloadJson = buffer.Bytes()
//	notification, err := NewNotification(GetIP(r), time.Now(), payloadJson)
//	err = notification.SaveToDb()
//	if handleRequestProcessingError(w, err, "Unable to save notification record to DB", 500) {
//		return
//	}
//	inbox.Add(*notification)
//	var page = NewPage()
//	page.Params["notificationUrl"] = fmt.Sprint(notification.ID)
//	page.Title = "Notification Response"
//	w.Header().Set("Location", notification.Url())
//	err = pageRender.HTML(w, http.StatusCreated, "post_success", page)
//	if handleRequestProcessingError(w, err, "Unable to process request", 500) {
//		return
//	}
//}

func InboxPost(w http.ResponseWriter, r *http.Request) {
	var err error
	payloadJson, err := ioutil.ReadAll(r.Body)
	if handleRequestProcessingError(w, err, "Unable to read posted content", 400) {
		return
	}
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, payloadJson)
	if handleRequestProcessingError(w, err, "Unable to parse posted content (must be JSON-LD)", 400) {
		return
	}
	payloadJson = buffer.Bytes()
	notification, err := NewNotification(GetIP(r), time.Now(), payloadJson)
	//err = notification.SaveToDb()
	if handleRequestProcessingError(w, err, "Unable to save notification record to DB", 500) {
		return
	}
	inbox.Add(*notification)
	var page = NewPage()
	page.Params["notificationUrl"] = fmt.Sprint(notification.ID)
	page.Title = "Notification Response"
	w.Header().Set("Location", notification.Url())
	err = pageRender.HTML(w, http.StatusCreated, "post_success", page)
	if handleRequestProcessingError(w, err, "Unable to process request", 500) {
		return
	}
}

func handleRequestProcessingError(w http.ResponseWriter, err error, message string, code int) bool {
	if err != nil {
		zapLogger.Error(message + err.Error())
		http.Error(w, message, code)
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
