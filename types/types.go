package types

import (
	"time"
	"net/url"
	"strings"
)

const (
	IDEA_COMPLETE    = 1 << iota
	USER_CREATED
	COMMENT_COMPLETE
)

type User struct {
	Id            string
	Username      string
	About         string
	Email         string
	AvatarUrl     string
	LinkedInUrl   string
	WebsiteUrl    string
	TwitterUrl    string
	MastodonUrl   string
	GithubUrl     string
	GitalUrl      string
	HackerNewsUrl string
	MediumUrl     string
	CreatedAt     time.Time
}

type UserValidUrls struct {
	Valid      bool
	Avatar     bool
	LinkedIn   bool
	Website    bool
	Twitter    bool
	Github     bool
	HackerNews bool
	Medium     bool
}

func (u User) GetState() (state uint16) {
	if u.Id != "" {
		state |= USER_CREATED
	}

	return
}

func (u User) GetValidUrls() (validUrls UserValidUrls) {
	_, err := url.ParseRequestURI("https://" + u.AvatarUrl)
	if u.AvatarUrl == "" || err == nil {
		validUrls.Avatar = true
	}

	_, err = url.ParseRequestURI("https://" + u.LinkedInUrl)
	if u.LinkedInUrl == "" || err == nil {
		validUrls.LinkedIn = true
	}

	_, err = url.ParseRequestURI("https://" + u.WebsiteUrl)
	if u.WebsiteUrl == "" || err == nil {
		validUrls.Website = true
	}

	_, err = url.ParseRequestURI("https://" + u.TwitterUrl)
	if u.TwitterUrl == "" || err == nil {
		validUrls.Twitter = true
	}

	_, err = url.ParseRequestURI("https://" + u.GithubUrl)
	if u.GithubUrl == "" || err == nil {
		validUrls.Github = true
	}

	_, err = url.ParseRequestURI("https://" + u.HackerNewsUrl)
	if u.HackerNewsUrl == "" || err == nil {
		validUrls.HackerNews = true
	}

	_, err = url.ParseRequestURI("https://" + u.MediumUrl)
	if u.MediumUrl == "" || err == nil {
		validUrls.Medium = true
	}

	if !validUrls.Avatar ||
		!validUrls.LinkedIn ||
		!validUrls.Website ||
		!validUrls.Twitter ||
		!validUrls.Github ||
		!validUrls.HackerNews ||
		!validUrls.Medium {
		validUrls.Valid = false
	} else {
		validUrls.Valid = true
	}

	return
}

type Idea struct {
	Id                  int64
	Upvoted             bool
	NUpvotes            int64
	NComments           int64
	AgeLabel            string
	Title               string
	Author              string
	Slug                string
	Badges              string
	DescriptionMarkdown string
	DescriptionHtml     string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	CreatedBy           string
}

func (i Idea) GetState() (state uint16) {
	if i.Title != "" && i.DescriptionMarkdown != "" {
		state |= IDEA_COMPLETE
	}

	return
}

func (i Idea) GetBadges() (badges []Badge) {
	for _, s := range strings.Split(i.Badges, ",") {
		badges = append(badges, NamedBadges[s])
	}

	return
}

type Upvote struct {
	UserId string
	IdeaId int64
}

type UpvoteCount struct {
	IdeaId int64
	Count  int64
}

type Comment struct {
	Id        int64
	IdeaId    int64
	UserId    string
	Username  string
	ParentId  int64
	Children  []*Comment
	CreatedAt time.Time
	AgeLabel  string
	Comment   string
	Level     int64
}

type SortCommentsByCreatedAt []*Comment

func (c SortCommentsByCreatedAt) Len() int           { return len(c) }
func (c SortCommentsByCreatedAt) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c SortCommentsByCreatedAt) Less(i, j int) bool { return c[i].CreatedAt.Before(c[j].CreatedAt) }

func (u Comment) GetState() (state uint16) {
	if u.IdeaId != 0 &&
		u.UserId != "" &&
		u.Username != "" &&
		u.Comment != "" &&
		!u.CreatedAt.IsZero() {
		state |= COMMENT_COMPLETE
	}

	return
}

type Visit struct {
	City          string
	CountryCode   string
	ContinentCode string
	Path          string
	VisitorId     string
	CreatedAt     time.Time
}

type VisitsCount struct {
	Count int64
	Date  string
}

type VisitsCounts []VisitsCount

func (vcs VisitsCounts) ToSeries() (series VisitsCountSeries) {
	series.Count = []int64{}
	series.Date = []string{}

	for _, vc := range vcs {
		series.Count = append(series.Count, vc.Count)
		series.Date = append(series.Date, vc.Date)
	}

	if len(series.Date) > 0 {
		series.Date[len(series.Date)-1] = "today"
	}

	if len(series.Date) == 1 {
		date := []string{}
		date = append(date, time.Now().AddDate(0, 0, -1).Format("01-02"))
		date = append(date, series.Date[0])
		series.Date = date

		count := []int64{}
		count = append(count, 0)
		count = append(count, series.Count[0])
		series.Count = count
	}

	return
}

type VisitsCountSeries struct {
	Count []int64  `json:"count"`
	Date  []string `json:"date"`
}
