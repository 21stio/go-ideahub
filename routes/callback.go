package routes

import (
	"context"
	_ "crypto/sha512"
	"encoding/json"
	"net/http"
	"github.com/gorilla/sessions"
	sess "github.com/21stio/go-ideahub/session"
	"github.com/21stio/go-ideahub/types"
	"github.com/21stio/go-ideahub/utils"
	"github.com/21stio/go-ideahub/queries"
	"time"
	"errors"
)

func GetCallBackHandler(store sessions.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			conf := utils.GetOAuthConfig()

			state := r.URL.Query().Get("state")
			session, err := store.Get(r, "state")
			if err != nil {
				return
			}

			if state != session.Values["state"] {
				err = errors.New("Invalid state parameter")

				return
			}

			code := r.URL.Query().Get("code")

			token, err := conf.Exchange(context.TODO(), code)
			if err != nil {
				return
			}

			domain := utils.GetEnv("AUTH0_DOMAIN")
			client := conf.Client(context.TODO(), token)
			resp, err := client.Get("https://" + domain + "/userinfo")
			if err != nil {
				return
			}

			defer resp.Body.Close()

			var profile map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&profile)
			if err != nil {
				return
			}

			session, err = store.Get(r, sess.AUTH)
			if err != nil {
				return
			}

			session.Values["id_token"] = token.Extra("id_token")
			session.Values["access_token"] = token.AccessToken
			session.Values["profile"] = profile
			err = session.Save(r, w)
			if err != nil {
				return
			}

			profileData := utils.GetProfileData(*session)

			user, err := queries.SelectUserById(profileData["sub"])
			if err != nil {
				return
			}

			if user.GetState()&types.USER_CREATED == 0 {
				user.Id = profileData["sub"]
				user.Username = profileData["nickname"]
				user.AvatarUrl = profileData["picture"]
				user.CreatedAt = time.Now()

				err = queries.InsertUser(user)
				if err != nil {
					return
				}
			}

			err = sess.SetUpvotesInSession(store, r, w, user.Id)
			if err != nil {
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
