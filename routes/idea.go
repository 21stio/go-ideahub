package routes

import (
	"github.com/21stio/go-ideahub/routes/templates"
	"github.com/21stio/go-ideahub/types"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"github.com/21stio/go-ideahub/queries"
)

type IdeaData struct {
	templates.BaseData
	Idea              types.Idea
	Comments          []*types.Comment
	Author            types.User
	DescriptionHtml   template.HTML
	VisitsCountSeries types.VisitsCountSeries
}

func GetIdeaHandler(tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			idea, err := queries.SelectIdeaBySlug(mux.Vars(r)["slug"])
			if err != nil {
				return
			}

			author, err := queries.SelectUserById(idea.CreatedBy)
			if err != nil {
				return
			}

			comments, err := queries.SelectComments(idea.Id)
			if err != nil {
				return
			}

			visitCountsSeries, err := queries.GetVisitsCountSeriesByIdeaSlug(idea.Slug)
			if err != nil {
				return
			}

			comments = addAgeLabelToComments(comments)

			comments = addHierachyToComments(comments)

			baseData = addUserToBaseData(baseData, author)

			ideaData := IdeaData{}
			ideaData.BaseData = baseData
			ideaData.Idea = idea
			ideaData.VisitsCountSeries = visitCountsSeries
			ideaData.TitleUrl = idea.Title
			ideaData.Comments = comments
			ideaData.Author = author
			ideaData.DescriptionHtml = template.HTML(idea.DescriptionHtml)

			err = tpl.ExecuteTemplate(w, "idea.html", ideaData)
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
