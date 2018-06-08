package middleware

import (
	"net/http"
	"github.com/gorilla/sessions"
	"context"
	sess "github.com/21stio/go-ideahub/session"
	"strings"
)

func InjectSession(store sessions.Store) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		err := func() (err error) {
			if strings.Contains(r.URL.Path, "public") {
				return
			}

			baseData, err := sess.GetBaseData(store, r)
			if err != nil {
				return
			}

			err = sess.UpdatePath(store, r, w)
			if err != nil {
				return
			}

			baseData.Path = sess.GetPath(store, r)

			baseData.Upvoted, err = sess.GetUpvotedFromSession(store, r)
			if err != nil {
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "baseData", baseData))

			return
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next(w, r)
	}
}
