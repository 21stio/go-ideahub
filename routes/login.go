package routes

import (
	"golang.org/x/oauth2"
	"net/http"
	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/sessions"
	sess "github.com/21stio/go-ideahub/session"
	"github.com/21stio/go-ideahub/utils"
	"os"
)

func GetLoginHandler(store sessions.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			domain := utils.GetEnv("AUTH0_DOMAIN")
			aud := os.Getenv("AUTH0_AUDIENCE")

			conf := utils.GetOAuthConfig()

			if aud == "" {
				aud = "https://" + domain + "/userinfo"
			}

			b := make([]byte, 32)
			rand.Read(b)
			state := base64.StdEncoding.EncodeToString(b)

			err = sess.SetState(store, r, w, state)
			if err != nil {
				return
			}

			audience := oauth2.SetAuthURLParam("audience", aud)
			url := conf.AuthCodeURL(state, audience)

			http.Redirect(w, r, url, http.StatusTemporaryRedirect)

			return
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
