package main

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gorilla/sessions"
	casbinFSAdapter "github.com/naucon/casbin-fs-adapter"
	"go.uber.org/zap"
	"net/http"
)

//go:embed auth_rules/authz_model.conf auth_rules/authz_policy.csv
var authFS embed.FS

func InitialiseCasbinEnforcer() *casbin.Enforcer {
	model, _ := casbinFSAdapter.NewModel(authFS, "auth_rules/authz_model.conf")
	policies := casbinFSAdapter.NewAdapter(authFS, "auth_rules/authz_policy.csv")
	enforcer, _ := casbin.NewEnforcer(model, policies)
	_ = enforcer.LoadPolicy()
	return enforcer
}

// Authorizer this function is adapted from https://github.com/casbin/chi-authz
func Authorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			session, _ := getSession(w, r)
			var user string
			//var user interface{}
			if session != nil {
				user = fmt.Sprint(session.Values["uid"])
			}
			method := r.Method
			path := r.URL.Path
			zapLogger.Debug("Authorising:", zap.Any("user", user), zap.Any("path", path), zap.Any("method", method))
			enforced, err := e.Enforce(user, path, method)
			if err != nil {
				zapLogger.Error(err.Error())
			}
			if enforced == true {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, http.StatusText(403), 403)
			}
		}
		return http.HandlerFunc(fn)
	}
}

func isAdmin(s *sessions.Session) bool {
	for _, adminEmail := range config.AdminEmails {
		if adminEmail == s.Values["uid"] {
			return true
		}
	}
	return false
}
