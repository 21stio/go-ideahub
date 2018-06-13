package routes

import (
	"net/http"
	"html/template"
	"github.com/21stio/go-ideahub/routes/templates"
	"github.com/21stio/go-ideahub/queries"
	"github.com/21stio/go-ideahub/types"
	"net/url"
	"strings"
	log "github.com/sirupsen/logrus"
)

type MeData struct {
	templates.BaseData
	Ideas     []types.Idea
	Me        types.User
	ValidUrls types.UserValidUrls
	Method    string
}

func GetGetMeHandler(tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			me, err := queries.SelectUserById(baseData.UserId)
			if err != nil {
				return
			}

			baseData = addUserToBaseData(baseData, me)

			ideas, err := queries.SelectIdeasByAuthorId(baseData.UserId)
			if err != nil {
				return
			}

			ideas = enrichIdeas(ideas, baseData)

			meData := MeData{}
			meData.BaseData = baseData
			meData.Ideas = ideas
			meData.Me = me
			meData.HtmlTitle = me.Username
			meData.SubmitOn = true

			err = tpl.ExecuteTemplate(w, "me.html", meData)
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

func GetPostMeHandler(tpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() (err error) {
			baseData, _ := r.Context().Value("baseData").(templates.BaseData)

			r.ParseForm()
			me := getMe(r.Form)
			me.Id = baseData.UserId
			validUrls := me.GetValidUrls()

			if validUrls.Valid {
				err = queries.UpdateUserById(me)
				if err != nil {
					return
				}

				http.Redirect(w, r, HOME, http.StatusTemporaryRedirect)

				return
			}

			ideas, err := queries.SelectIdeasByAuthorId(baseData.UserId)
			if err != nil {
				return
			}

			ideas = enrichIdeas(ideas, baseData)

			me.Username = baseData.Username
			baseData = addUserToBaseData(baseData, me)

			meData := MeData{}
			meData.BaseData = baseData
			meData.Ideas = ideas
			meData.Me = me
			meData.ValidUrls = validUrls
			meData.HtmlTitle = me.Username
			meData.Method = http.MethodPost

			err = tpl.ExecuteTemplate(w, "me.html", meData)
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

func removeHttpHttps(s string) (string) {
	return strings.Replace(strings.Replace(s, "http://", "", 1), "https://", "", 1)
}

func getMe(form url.Values) (user types.User) {
	user.About = form.Get("about")
	user.Email = form.Get("email")
	user.AvatarUrl = removeHttpHttps(form.Get("avatar_url"))
	user.WebsiteUrl = removeHttpHttps(form.Get("website_url"))
	user.HackerNewsUrl = removeHttpHttps(form.Get("hackernews_url"))
	user.TwitterUrl = removeHttpHttps(form.Get("twitter_url"))
	user.GithubUrl = removeHttpHttps(form.Get("github_url"))
	user.MediumUrl = removeHttpHttps(form.Get("medium_url"))
	user.LinkedInUrl = removeHttpHttps(form.Get("linkedin_url"))

	return
}
