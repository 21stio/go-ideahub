package routes

import (
	"net/http"
	"github.com/21stio/go-ideahub/routes/templates"
	"github.com/21stio/go-ideahub/queries"
	"github.com/gorilla/mux"
	"github.com/21stio/go-ideahub/types"
	"strconv"
	"strings"
	sess "github.com/21stio/go-ideahub/session"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

func GetUpvotedHandler(store sessions.Store) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			vars := mux.Vars(r)

			ideaId, err := strconv.ParseInt(vars["idea_id"], 10, 64)
			if err != nil {
				return
			}

			upvote := types.Upvote{}
			upvote.IdeaId = ideaId
			upvote.UserId = baseData.UserId

			err = queries.InsertUpvote(upvote)
			if err != nil && isDuplicateKeyError(err) {
				err = nil
			}
			if err != nil {
				return
			}

			err = sess.AddUpvotedToSession(store, r, w, ideaId)
			if err != nil {
				return
			}

			http.Redirect(w, r, baseData.Path.Previous, http.StatusTemporaryRedirect)

			return
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
		}
	}
}

func isDuplicateKeyError(err error) (bool) {
	return strings.Contains(err.Error(), "duplicate key value")
}
