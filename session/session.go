package session

import (
	"github.com/gorilla/sessions"
	"github.com/davecgh/go-spew/spew"
	"github.com/21stio/go-ideahub/queries"
	"errors"
	"net/http"
	"github.com/21stio/go-ideahub/routes/templates"
)

const (
	AUTH    = "auth"
	PATH    = "path"
	UPVOTES = "upvotes"
	STATE   = "state"
)

func GetBaseData(store sessions.Store, r *http.Request) (baseData templates.BaseData, err error) {
	session, err := store.Get(r, AUTH)
	if err != nil {
		return
	}

	baseData = templates.BaseData{}
	baseData.HtmlTitle = "sharing ideas"
	baseData.Title = "TMRW"
	baseData.SubTitle = "sharing ideas"
	baseData.TitleUrl = "/"
	baseData.LogoUrl = "/public/littleguy.png"

	_, ok := session.Values["id_token"]
	baseData.IsLoggedIn = ok

	if ok {
		profile := getProfileData(session)

		baseData.Username = getStringFromProfile(profile, "nickname")
		baseData.UserId = getStringFromProfile(profile, "sub")
	}

	return
}

func UpdatePath(store sessions.Store, r *http.Request, w http.ResponseWriter) (err error) {
	session, err := store.Get(r, PATH)
	if err != nil {
		return
	}

	session.Values["previous"] = session.Values["current"]
	session.Values["current"] = r.URL.Path

	err = session.Save(r, w)
	if err != nil {
		return
	}

	return
}

func GetPath(store sessions.Store, r *http.Request) (path templates.Path) {
	session, err := store.Get(r, PATH)
	if err != nil {
		return
	}

	path.Current = getStringFromSessionValues(session.Values, "current")
	path.Previous = getStringFromSessionValues(session.Values, "previous")

	return
}

func SetState(store sessions.Store, r *http.Request, w http.ResponseWriter, state string) (err error) {
	session, err := store.Get(r, STATE)
	if err != nil {
		return
	}

	session.Values[STATE] = state
	err = session.Save(r, w)
	if err != nil {
		return
	}

	return
}

func SetUpvotesInSession(store sessions.Store, r *http.Request, w http.ResponseWriter, userId string) (err error) {
	session, err := store.Get(r, UPVOTES)
	if err != nil {
		return
	}

	upvotes, err := queries.SelectUpvotesByUserId(userId)
	if err != nil {
		return
	}

	ids := []int64{}
	for _, upvote := range upvotes {
		ids = append(ids, upvote.IdeaId)
	}

	session.Values[UPVOTES] = ids

	err = session.Save(r, w)
	if err != nil {
		return
	}

	return
}

func GetUpvotedFromSession(store sessions.Store, r *http.Request) (upvoted map[int64]bool, err error) {
	session, err := store.Get(r, UPVOTES)
	if err != nil {
		return
	}

	ids, err := getUpvoteIds(session)
	if err != nil {
		return
	}

	upvoted = map[int64]bool{}

	for _, id := range ids {
		upvoted[id] = true
	}

	return
}

func AddUpvotedToSession(store sessions.Store, r *http.Request, w http.ResponseWriter, ideaId int64) (err error) {
	session, err := store.Get(r, UPVOTES)
	if err != nil {
		return
	}

	ids, err := getUpvoteIds(session)
	if err != nil {
		return
	}

	ids = append(ids, ideaId)

	session.Values[UPVOTES] = ids

	err = session.Save(r, w)
	if err != nil {
		return
	}

	return
}

func GetState(store sessions.Store, r *http.Request) (state interface{}, err error) {
	session, err := store.Get(r, STATE)
	if err != nil {
		return
	}

	state = session.Values["state"]

	return
}

func Clear(store sessions.Store, r *http.Request, w http.ResponseWriter) (err error) {
	var session *sessions.Session

	for _, s := range []string{AUTH, PATH, UPVOTES, STATE} {
		session, err = store.Get(r, s)
		if err != nil {
			return
		}

		session.Values = map[interface{}]interface{}{}

		err = session.Save(r, w)
		if err != nil {
			return
		}
	}

	return
}

func getUpvoteIds(session *sessions.Session) (ids []int64, err error) {
	if session.Values[UPVOTES] == nil {
		return
	}

	ids, ok := session.Values[UPVOTES].([]int64)
	if !ok {
		err = errors.New("not []int64")

		return
	}

	return
}

func getProfileData(session *sessions.Session) (profile map[string]interface{}) {
	profile, ok := session.Values["profile"].(map[string]interface{})
	if !ok {
		spew.Dump("count not retrieve profile data")
	}

	return
}

func getStringFromProfile(profile map[string]interface{}, key string) (string) {
	value, ok := profile[key].(string)
	if !ok {
		spew.Dump("could not get string from profile key:" + key)
	}

	return value
}

func getStringFromSessionValues(values map[interface{}]interface{}, key string) (string) {
	value, ok := values[key].(string)
	if !ok {
		spew.Dump("could not get string from session values key:" + key)
	}

	return value
}
