package utils

import (
	"github.com/gorilla/sessions"
	"github.com/davecgh/go-spew/spew"
	"os"
	"golang.org/x/oauth2"
	"log"
	"fmt"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
}

func GetProfileData(session sessions.Session) (data map[string]string) {
	data = map[string]string{}

	profile, ok := session.Values["profile"].(map[string]interface{})
	if !ok {
		spew.Dump("a")
	}

	a := func(key string, d map[string]interface{}) (string) {
		value, ok := profile[key].(string)
		if !ok {
			spew.Dump("b")
		}

		return value
	}

	keys := []string{"name", "nickname", "sub", "updated_at", "picture"}

	for _, key := range keys {
		data[key] = a(key, profile)
	}

	return
}

func GetOAuthConfig() (conf *oauth2.Config) {
	domain := GetEnv("AUTH0_DOMAIN")

	conf = &oauth2.Config{
		ClientID:     GetEnv("AUTH0_CLIENT_ID"),
		ClientSecret: GetEnv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  GetEnv("AUTH0_CALLBACK_URL"),
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	return
}

func GetUniqueStrings(s []string) []string {
	u := make([]string, 0, len(s))
	m := make(map[string]bool)

	for _, val := range s {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func GetEnv(key string) (string) {
	s := os.Getenv(key)
	if s == "" {
		log.Fatal(fmt.Sprintf("env key:%v not set", key))
	}

	return s
}
