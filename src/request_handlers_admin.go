package main

import (
	"github.com/johnsto/go-passwordless"
	"net/http"
)

func AdminPageGet(w http.ResponseWriter, r *http.Request) {
	var page = NewPage()
	page.Title = "Admin"
	pageBytes, _ := GetPageBodyAsByteSliceFromFs("admin.md")
	html, _ := GetHTMLFromMarkdown(pageBytes)
	page.Body = html
	pageRender.HTML(w, http.StatusOK, "page", page)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if session, err := getSession(w, r); err == nil {
		if isSignedIn(session) {
			session.AddFlash("already_signed_in")
			session.Save(r, w)
			redirect(w, r, "/", site.BaseUrl)
			return
		}
	}
	page := NewPage()
	page.Title = "Login"
	page.Params["Next"] = r.FormValue("next")
	pageRender.HTML(w, http.StatusOK, "login", page)
}

func TokenCheck(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(w, r)
	if err != nil {
		zapLogger.Error(err.Error())
		return
	}
	if isSignedIn(session) {
		session.AddFlash("already_signed_in")
		session.Save(r, w)
		redirect(w, r, r.FormValue("next"), site.BaseUrl)
		return
	}
	ctx := passwordless.SetContext(nil, w, r)
	strategy := "email"
	recipient := r.FormValue("recipient")
	token := r.FormValue("token")
	uid := recipient
	tokenError := ""
	if token == "" {
		err := PW.RequestToken(ctx, strategy, uid, recipient)
		if err != nil {
			zapLogger.Error(err.Error())
			http.Error(w, "Internal error"+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// User has provided a token, verify it against provided uid.
		valid, err := PW.VerifyToken(ctx, uid, token)
		if valid {
			// User provided a valid token! We can safely use the uid as it
			// is validated alongside the token.
			session.Values["uid"] = uid
			session.AddFlash("signed_in")
			session.Save(r, w)
			redirect(w, r, r.FormValue("next"), site.BaseUrl)
			return
		}
		if err == passwordless.ErrTokenNotFound {
			// Token not found, maybe it was a previous one or expired. Either
			// way, the user will need to attempt sign-in again.
			session.AddFlash("token_not_found")
			session.Save(r, w)
			http.Redirect(w, r, "/account/signin", http.StatusTemporaryRedirect)
			return
		} else if err != nil {
			// Some other unexpected error occurred.
			zapLogger.Error(err.Error())
			http.Error(w, "Failed verifying token"+err.Error(), http.StatusInternalServerError)
			return
		} else {
			// User entered bad token. Set token error string then fall
			// through to template.
			w.WriteHeader(http.StatusForbidden)
			tokenError = "The entered token/PIN was incorrect."
		}
	}
	// If we've got to this point, the user is being prompted to enter a
	// valid token value.
	page := NewPage()
	page.Title = "Enter Token"
	page.Params["Next"] = r.FormValue("next")
	page.Params["Recipient"] = recipient
	page.Params["UserID"] = uid
	page.Params["TokenError"] = tokenError
	pageRender.HTML(w, http.StatusOK, "token", page)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(w, r)
	if err != nil {
		return
	}
	// Remove secure session cookie
	delete(session.Values, "uid")
	session.AddFlash("signed_out")
	session.Save(r, w)
	//redirect(w, r, r.FormValue("next"), site.BaseUrl)
	http.Redirect(w, r, site.BaseUrl, http.StatusFound)
}
