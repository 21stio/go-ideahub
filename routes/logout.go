package routes

import (
	"net/http"
	"net/url"
	"github.com/gorilla/sessions"
	"github.com/21stio/go-ideahub/session"
	"github.com/21stio/go-ideahub/utils"
	log "github.com/sirupsen/logrus"
)

func GetLogoutHandler(store sessions.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			domain := utils.GetEnv("AUTH0_DOMAIN")

			var Url *url.URL
			Url, err = url.Parse("https://" + domain)
			if err != nil {
				return
			}

			err = session.Clear(store, r, w)
			if err != nil {
				return
			}

			Url.Path += "/v2/logout"
			parameters := url.Values{}
			parameters.Add("returnTo", utils.GetEnv("LOGOUT_RETURN_TO"))
			parameters.Add("client_id", utils.GetEnv("AUTH0_CLIENT_ID"))
			Url.RawQuery = parameters.Encode()

			http.Redirect(w, r, Url.String(), http.StatusTemporaryRedirect)

			return
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
		}
	}
}
