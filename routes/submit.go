package routes

import (
	"net/http"
	"html/template"
	"github.com/21stio/go-ideahub/routes/templates"
	"github.com/21stio/go-ideahub/types"
	"gopkg.in/russross/blackfriday.v2"
	"github.com/microcosm-cc/bluemonday"
	"net/url"
	"github.com/21stio/go-ideahub/queries"
	"time"
	"github.com/gosimple/slug"
	"math/rand"
	"strings"
	"github.com/davecgh/go-spew/spew"
)

type SubmitData struct {
	templates.BaseData
	Idea   types.Idea
	Method string
	Badges []types.Badge
}

func GetGetSubmitHandler(tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			me, err := queries.SelectUserById(baseData.UserId)
			if err != nil {
				return
			}

			baseData = addUserToBaseData(baseData, me)

			submitData := SubmitData{}
			submitData.BaseData = baseData
			submitData.Badges = types.Badges

			err = tpl.ExecuteTemplate(w, "submit.html", submitData)
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

func GetPostSubmitHandler(tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			r.ParseForm()

			spew.Dump(r.Form)

			idea := getIdea(r.Form)
			idea.CreatedBy = baseData.UserId

			if idea.GetState()&types.IDEA_COMPLETE != 0 {
				err = queries.InsertIdea(idea)
				if err != nil {
					return
				}

				http.Redirect(w, r, HOME, http.StatusTemporaryRedirect)

				return
			}

			me, err := queries.SelectUserById(baseData.UserId)
			if err != nil {
				return
			}

			baseData = addUserToBaseData(baseData, me)

			submitData := SubmitData{}
			submitData.BaseData = baseData
			submitData.Idea = idea
			submitData.Badges = types.Badges
			submitData.Method = http.MethodPost

			err = tpl.ExecuteTemplate(w, "submit.html", submitData)
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

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getIdea(form url.Values) (idea types.Idea) {
	idea.Title = form.Get("title")
	idea.Slug = randSeq(5) + "-" + slug.Make(form.Get("title"))
	idea.Badges = strings.Join(form["badges"], ",")
	idea.DescriptionMarkdown = form.Get("description")
	idea.CreatedAt = time.Now().UTC()

	unsafe := blackfriday.Run([]byte(idea.DescriptionMarkdown))
	idea.DescriptionHtml = string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))

	return
}
