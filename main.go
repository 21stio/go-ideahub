package main

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/codegangsta/negroni"
	"github.com/21stio/go-ideahub/routes"
	"github.com/21stio/go-ideahub/middleware"
	"os"
	"html/template"
	"github.com/21stio/go-ideahub/queries"
	"database/sql"
	_ "github.com/lib/pq"
	"reflect"
	"net/url"
	"github.com/oschwald/geoip2-golang"
	"encoding/json"
	"fmt"
)

var (
	store  *sessions.CookieStore
	tpl    *template.Template
	reader *geoip2.Reader
)

func init() {
	var err error
	queries.Db, err = sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	var fns = template.FuncMap{
		"sub": func(x int64, y int64) int64 {
			return x - y
		},
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"queryEscape": url.QueryEscape,
		"marshal": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tpl, err = template.New("").Funcs(fns).ParseGlob(cwd + "/routes/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	reader, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	gob.Register(map[string]interface{}{})

	http.Handle("/", getMiddleware())

	addr := os.Getenv("ADDR")
	fmt.Printf("Server listening on http://%v/", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getMiddleware() (n *negroni.Negroni) {
	root := mux.NewRouter()
	root.HandleFunc("/", routes.GetHomeHandler(tpl))
	root.HandleFunc(routes.IDEA+"/{slug}", routes.GetIdeaHandler(tpl))
	root.HandleFunc(routes.USER+"/{username}", routes.GetUserHandler(tpl))
	root.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

	auth := root.PathPrefix(routes.AUTH).Subrouter()
	auth.HandleFunc(routes.LOGIN, routes.GetLoginHandler(store))
	auth.HandleFunc(routes.LOGIN_CALLBACK, routes.GetCallBackHandler(store))
	auth.HandleFunc(routes.LOGOUT, routes.GetLogoutHandler(store))

	authedBase := mux.NewRouter()
	root.PathPrefix(routes.AUTHENTICATED).Handler(negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(negroni.HandlerFunc(middleware.GetIsAuthenticated(store))),
		negroni.Wrap(authedBase),
	))

	authed := authedBase.PathPrefix(routes.AUTHENTICATED).Subrouter()
	authed.HandleFunc(routes.ME, routes.GetGetMeHandler(tpl)).Methods(http.MethodGet)
	authed.HandleFunc(routes.ME, routes.GetPostMeHandler(tpl)).Methods(http.MethodPost)
	authed.HandleFunc(routes.SUBMIT, routes.GetGetSubmitHandler(tpl)).Methods(http.MethodGet)
	authed.HandleFunc(routes.SUBMIT, routes.GetPostSubmitHandler(tpl)).Methods(http.MethodPost)
	authed.HandleFunc(routes.UPVOTE+"/{idea_id}/{slug}", routes.GetUpvotedHandler(store))
	authed.HandleFunc(routes.COMMENT, routes.GetPostCommentHandler()).Methods(http.MethodPost)

	n = negroni.New(
		negroni.NewLogger(),
		negroni.HandlerFunc(middleware.InjectSession(store)),
		negroni.HandlerFunc(middleware.Tracking(reader)),
	)
	n.UseHandler(root)

	return
}
