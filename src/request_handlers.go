package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/alecthomas/chroma/formatters"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
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
	payloadJson, _ := ioutil.ReadAll(r.Body)
	buffer := new(bytes.Buffer)
	json.Compact(buffer, payloadJson)
	payloadJson = buffer.Bytes()
	notification := NewNotification(GetIP(r), time.Now(), payloadJson)
	notification.SaveNotificationToDb()
	inbox.Add(*notification)
	var page = NewPage()
	page.Params["notificationUrl"] = fmt.Sprint(notification.ID)
	page.Title = "Notification Response"
	w.Header().Set("Location", notification.Url())
	pageRender.HTML(w, http.StatusCreated, "post_success", page)
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
	id, _ := strconv.ParseUint(idString, 10, 32)
	notification := LoadNotificationFromDbById(uint(id))
	var page = NewNotificationPage(notification)
	page.Title = "NotificationRecord"
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
