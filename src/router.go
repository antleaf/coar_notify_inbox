package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func ConfigureRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	//r.Use(cors.Handler(cors.Options{
	//	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
	//	AllowedOrigins: []string{"https://*", "http://*"},
	//	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	//	AllowedMethods:   []string{"GET", "HEAD", "POST", "OPTIONS"},
	//	//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	//	ExposedHeaders:   []string{"Link"},
	//	AllowCredentials: false,
	//	MaxAge:           300, // Maximum value not ignored by any of major browsers
	//}))
	r.Handle("/assets/*", http.FileServer(http.FS(embeddedAssets)))
	r.Get("/", HomePageGet)
	r.Get("/inbox", InboxGet)
	r.Get("/inbox.json", InboxGetJSON)
	r.Post("/inbox", InboxPost)
	r.Get("/inbox/{id}", InboxNotificationGet)
	r.Options("/inbox", Options)
	return r
}
