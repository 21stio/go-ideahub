package routes

import (
	"github.com/21stio/go-ideahub/routes/templates"
	"github.com/21stio/go-ideahub/queries"
	"github.com/21stio/go-ideahub/types"
	"net/http"
	"net/url"
	"time"
	"strconv"
	log "github.com/sirupsen/logrus"
)

func GetPostCommentHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			r.ParseForm()
			comment, err := getComment(r.Form, baseData)
			if err != nil {
				return
			}

			if comment.GetState()&types.COMMENT_COMPLETE != 0 {
				err = queries.InsertComment(comment)
				if err != nil {
					return
				}
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

func getComment(form url.Values, data templates.BaseData) (comment types.Comment, err error) {
	parentId := int64(0)
	parentIdString := form.Get("parent_id")

	if parentIdString != "" {
		parentId, err = strconv.ParseInt(parentIdString, 10, 64)
		if err != nil {
			return
		}
	}

	ideaId := int64(0)
	ideaIdString := form.Get("idea_id")

	if ideaIdString != "" {
		ideaId, err = strconv.ParseInt(ideaIdString, 10, 64)
		if err != nil {
			return
		}
	}

	comment.IdeaId = ideaId
	comment.Comment = form.Get("comment")
	comment.UserId = data.UserId
	comment.Username = data.Username
	comment.ParentId = parentId
	comment.CreatedAt = time.Now().UTC()

	return
}
