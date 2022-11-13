package main

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/johnsto/go-passwordless"
	"gopkg.in/throttled/throttled.v2"
	"io"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"time"
)

var SessionStore sessions.Store
var PW *passwordless.Passwordless
var Limiter *throttled.HTTPRateLimiter

func initialisePasswordlessAuthentication() {
	SessionStore = sessions.NewCookieStore(config.CookieKey)
	tokStore := passwordless.NewMemStore()
	PW = passwordless.New(tokStore)
	smtpHostWithPort := os.Getenv("ANTLEAF_ROBOT_SMTP_HOST") + ":" + os.Getenv("ANTLEAF_ROBOT_SMTP_PORT")
	PW.SetTransport("email", passwordless.NewSMTPTransport(
		smtpHostWithPort,
		os.Getenv("Antleaf Robot"),
		smtp.PlainAuth(
			os.Getenv("ANTLEAF_ROBOT_SMTP_USERNAME"),
			os.Getenv("ANTLEAF_ROBOT_SMTP_USERNAME"),
			os.Getenv("ANTLEAF_ROBOT_SMTP_PASSWORD"),
			os.Getenv("ANTLEAF_ROBOT_SMTP_HOST")),
		emailWriter,
	), passwordless.NewCrockfordGenerator(10), 30*time.Minute)
}

func emailWriter(ctx context.Context, token, uid, recipient string, w io.Writer) error {
	e := &passwordless.Email{
		Subject: "Go-Passwordless signin",
		To:      recipient,
	}

	link := site.BaseUrl + "/account/token" +
		"?strategy=email&token=" + token + "&uid=" + uid

	// Ideally these would be populated from templates, but...
	text := "You (or someone who knows your email address) wants " +
		"to sign in to the Go-Passwordless website.\n\n" +
		"Your PIN is " + token + " - or use the following link: " +
		link + "\n\n" +
		"(If you were did not request or were not expecting this email, " +
		"you can safely ignore it.)"
	html := "<!doctype html><html><body>" +
		"<p>You (or someone who knows your email address) wants " +
		"to sign in to the Go-Passwordless website.</p>" +
		"<p>Your PIN is <b>" + token + "</b> - or <a href=\"" + link + "\">" +
		"click here</a> to sign in automatically.</p>" +
		"<p>(If you did not request or were not expecting this email, " +
		"you can safely ignore it.)</p></body></html>"

	// Add content types, from least- to most-preferable.
	e.AddBody("text/plain", text)
	e.AddBody("text/html", html)

	_, err := e.Write(w)

	return err
}

func getSession(w http.ResponseWriter, r *http.Request) (*sessions.Session, error) {
	session, err := SessionStore.Get(r, string(config.SessionKey))
	if err != nil && session == nil {
		session, err = SessionStore.New(r, string(config.SessionKey))
		session.Options.MaxAge = config.SessionExpiryIntervalSeconds //defaults to 86400
		if err != nil && session == nil {
			//writeError(w, r, session, http.StatusUnauthorized, Error{
			//	Name:        "Couldn't get session",
			//	Description: err.Error(),
			//	Error:       err,
			//})
			return nil, err
		}
	}
	return session, nil
}

func isSignedIn(s *sessions.Session) bool {
	return s != nil && s.Values["uid"] != nil
}

func redirect(w http.ResponseWriter, r *http.Request, next, base string) {
	if nextURL, err := url.Parse(next); err != nil {
		zapLogger.Error(err.Error())
		next = base
	} else if nextURL.IsAbs() && next[:len(base)] != base {
		zapLogger.Error(err.Error())
		next = base
	}
	http.Redirect(w, r, next, http.StatusFound)
}

// getTemplateContext returns a Context object containing the current user
// and other variables required by all templates.
