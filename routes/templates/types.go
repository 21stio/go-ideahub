package templates

type BaseData struct {
	IsLoggedIn bool
	IsHome     bool
	SubmitOn   bool
	Upvoted    map[int64]bool
	Username   string
	UserId     string
	HtmlTitle  string
	Title      string
	TitleUrl   string
	SubTitle   string
	LogoUrl    string
	Path       Path
}

type Path struct {
	Previous string
	Current  string
}
