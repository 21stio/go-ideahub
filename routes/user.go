package routes

import (
	"github.com/21stio/go-ideahub/routes/templates"
	"github.com/21stio/go-ideahub/queries"
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/21stio/go-ideahub/types"
)

type UserData struct {
	templates.BaseData
	User  types.User
	Ideas []types.Idea
	VisitsCountSeries types.VisitsCountSeries
}

func GetUserHandler(tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			user, err := queries.SelectUserByUsername(mux.Vars(r)["username"])
			if err != nil {
				return
			}

			ideas, err := queries.SelectIdeasByAuthorId(user.Id)
			if err != nil {
				return
			}

			visitCountsSeries, err := queries.GetVisitsCountSeriesByUsername(user.Username)
			if err != nil {
				return
			}

			ideas = enrichIdeas(ideas, baseData)

			baseData = addUserToBaseData(baseData, user)

			userData := UserData{}
			userData.BaseData = baseData
			userData.User = user
			userData.TitleUrl = user.Username
			userData.Ideas = ideas
			userData.VisitsCountSeries = visitCountsSeries

			err = tpl.ExecuteTemplate(w, "user.html", userData)
			if err != nil {
				return
			}

			return
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
