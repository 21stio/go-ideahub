package routes

import (
	"net/http"
	"html/template"
	"github.com/21stio/go-ideahub/routes/templates"
	"github.com/21stio/go-ideahub/queries"
	"github.com/21stio/go-ideahub/types"
	log "github.com/sirupsen/logrus"
)

type HomeData struct {
	templates.BaseData
	Ideas []types.Idea
}

func GetHomeHandler(tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			ideas, err := queries.SelectIdeas(0)
			if err != nil {
				return
			}

			ideas = enrichIdeas(ideas, baseData)

			homeData := HomeData{}
			homeData.BaseData = baseData
			homeData.Ideas = ideas
			homeData.SubmitOn = true
			homeData.IsHome = true

			err = tpl.ExecuteTemplate(w, "home.html", homeData)
			if err != nil {
				return
			}

			return
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
		}
	}
}
