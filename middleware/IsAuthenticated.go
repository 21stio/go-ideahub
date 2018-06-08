package middleware

import (
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/21stio/go-ideahub/routes/templates"
	"strings"
)

func GetIsAuthenticated(store sessions.Store) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if strings.Contains(r.URL.Path, "public") {
			return
		}

		baseData, _ := r.Context().Value("baseData").(templates.BaseData)

		if !baseData.IsLoggedIn {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		} else {
			next(w, r)
		}
	}
}
