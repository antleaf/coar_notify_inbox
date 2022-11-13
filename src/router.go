package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func ConfigureRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(Authorizer(InitialiseCasbinEnforcer()))
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Throttle(10))
	//TODO: figure out if it is possible to use this CORS module to add common HTTP headers to all HTTP Responses. Otherwise write a middleware handler to do this.
	//r.Use(cors.Handler(cors.Options{
	r.Handle("/assets/*", http.FileServer(http.FS(embeddedAssets)))
	r.Mount("/admin", adminRouter())
	r.Get("/", HomePageGet)
	r.Get("/inbox", InboxGet)
	r.Post("/inbox", InboxPost)
	r.Get("/inbox/{id}", InboxNotificationGet)
	r.Options("/inbox", Options)
	r.Get("/login", Login)
	r.Post("/token", TokenCheck)
	r.Get("/logout", Logout)
	return r
}

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//r.Use(middleware.BasicAuth("COAR Notify Inbox Admin Realm", map[string]string{
	//	"admin": os.Getenv("COAR_NOTIFY_INBOX_ADMIN_PASSWORD"),
	//}))
	r.Get("/", AdminPageGet)
	return r
}
