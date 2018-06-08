package queries

import (
	"github.com/21stio/go-ideahub/types"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/lib/pq"
	"fmt"
	"strings"
	"github.com/davecgh/go-spew/spew"
)

var Db *sql.DB

func InsertIdea(idea types.Idea) (err error) {
	spew.Dump(idea)

	_, err = Db.Exec("INSERT INTO ideas(title, slug, badges, description_markdown, description_html, created_at, updated_at, created_by) VALUES($1, $2, $3, $4, $5, $6, $7, $8)", idea.Title, idea.Slug, idea.Badges, idea.DescriptionMarkdown, idea.DescriptionHtml, idea.CreatedAt, idea.UpdatedAt, idea.CreatedBy)
	return
}

func SelectIdeas(page int64) (ideas []types.Idea, err error) {
	rows, err := Db.Query(`
WITH ideas_0 as (
  SELECT id, title, slug, badges, description_markdown, description_html, created_at, updated_at, created_by,
         (
              SELECT COUNT(*)
              FROM comments
              WHERE comments.idea_id=ideas.id
         ) AS n_comments,
         (
              SELECT COUNT(*)
              FROM upvotes
              WHERE upvotes.idea_id=ideas.id
         ) AS n_upvotes
  FROM ideas
)

SELECT *
FROM ideas_0
ORDER BY n_upvotes/(EXTRACT(EPOCH FROM current_timestamp-created_at)/3600)^1.8 desc
LIMIT 30`)
	if err != nil {
		return
	}

	for rows.Next() {
		idea := types.Idea{}
		err = rows.Scan(&idea.Id, &idea.Title, &idea.Slug, &idea.Badges, &idea.DescriptionMarkdown, &idea.DescriptionHtml, &idea.CreatedAt, &idea.UpdatedAt, &idea.CreatedBy, &idea.NComments, &idea.NUpvotes)
		if err != nil {
			return
		}
		ideas = append(ideas, idea)
	}

	return
}

func SelectIdeasByAuthorId(authorId string) (ideas []types.Idea, err error) {
	rows, err := Db.Query(`
WITH ideas_0 as (
  SELECT id, title, slug, badges, description_markdown, description_html, created_at, updated_at, created_by,
         (
              SELECT COUNT(*)
              FROM comments
              WHERE comments.idea_id=ideas.id
         ) AS n_comments,
         (
              SELECT COUNT(*)
              FROM upvotes
              WHERE upvotes.idea_id=ideas.id
         ) AS n_upvotes
  FROM ideas
  WHERE created_by=$1
)

SELECT *
FROM ideas_0
ORDER BY n_upvotes/(EXTRACT(EPOCH FROM current_timestamp-created_at)/3600)^1.8 desc`, authorId)
	if err != nil {
		return
	}

	for rows.Next() {
		idea := types.Idea{}
		err = rows.Scan(&idea.Id, &idea.Title, &idea.Slug, &idea.Badges, &idea.DescriptionMarkdown, &idea.DescriptionHtml, &idea.CreatedAt, &idea.UpdatedAt, &idea.CreatedBy, &idea.NComments, &idea.NUpvotes)
		if err != nil {
			return
		}
		ideas = append(ideas, idea)
	}

	return
}

func SelectIdeaBySlug(slug string) (idea types.Idea, err error) {
	rows, err := Db.Query("SELECT id, title, slug, badges, description_markdown, description_html, created_at, updated_at, created_by FROM ideas WHERE slug=$1", slug)
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&idea.Id, &idea.Title, &idea.Slug, &idea.Badges, &idea.DescriptionMarkdown, &idea.DescriptionHtml, &idea.CreatedAt, &idea.UpdatedAt, &idea.CreatedBy)
		if err != nil {
			return
		}
	}

	return
}

func SelectUserById(id string) (user types.User, err error) {
	rows, err := Db.Query("SELECT id, username, about, email, avatar_url, linkedin_url, website_url, twitter_url, github_url, hackernews_url, medium_url, created_at FROM users WHERE id=$1", id)
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.About, &user.Email, &user.AvatarUrl, &user.LinkedInUrl, &user.WebsiteUrl, &user.TwitterUrl, &user.GithubUrl, &user.HackerNewsUrl, &user.MediumUrl, &user.CreatedAt)
		if err != nil {
			return
		}
	}

	return
}

func SelectUsersByIds(ids []string) (users []types.User, err error) {
	rows, err := Db.Query("SELECT id, username, about, email, avatar_url, linkedin_url, website_url, twitter_url, github_url, hackernews_url, medium_url, created_at FROM users WHERE id=ANY($1)", pq.Array(ids))
	if err != nil {
		return
	}

	for rows.Next() {
		user := types.User{}
		err = rows.Scan(&user.Id, &user.Username, &user.About, &user.Email, &user.AvatarUrl, &user.LinkedInUrl, &user.WebsiteUrl, &user.TwitterUrl, &user.GithubUrl, &user.HackerNewsUrl, &user.MediumUrl, &user.CreatedAt)
		if err != nil {
			return
		}
		users = append(users, user)
	}

	return
}

func SelectUserByUsername(username string) (user types.User, err error) {
	rows, err := Db.Query("SELECT id, username, about, email, avatar_url, linkedin_url, website_url, twitter_url, github_url, hackernews_url, medium_url, created_at FROM users WHERE username=$1", username)
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.About, &user.Email, &user.AvatarUrl, &user.LinkedInUrl, &user.WebsiteUrl, &user.TwitterUrl, &user.GithubUrl, &user.HackerNewsUrl, &user.MediumUrl, &user.CreatedAt)
		if err != nil {
			return
		}
	}

	return
}

func UpdateUserById(user types.User) (err error) {
	_, err = Db.Exec("UPDATE users SET about=$2, email=$3, avatar_url=$4, linkedin_url=$5, website_url=$6, twitter_url=$7, github_url=$8, hackernews_url=$9, medium_url=$10 WHERE id=$1", user.Id, user.About, user.Email, user.AvatarUrl, user.LinkedInUrl, user.WebsiteUrl, user.TwitterUrl, user.GithubUrl, user.HackerNewsUrl, user.MediumUrl)
	if err != nil {
		return
	}

	return
}

func InsertUser(user types.User) (err error) {
	_, err = Db.Exec("INSERT INTO users(id, username, about, email, avatar_url, linkedin_url, website_url, twitter_url, github_url, hackernews_url, medium_url, created_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", user.Id, user.Username, user.About, user.Email, user.AvatarUrl, user.LinkedInUrl, user.WebsiteUrl, user.TwitterUrl, user.GithubUrl, user.HackerNewsUrl, user.MediumUrl, user.CreatedAt)

	return
}

func SelectUpvotesByUserId(userId string) (upvotes []types.Upvote, err error) {
	rows, err := Db.Query("SELECT idea_id FROM upvotes WHERE user_id=$1", userId)

	for rows.Next() {
		upvote := types.Upvote{}
		upvote.UserId = userId

		err = rows.Scan(&upvote.IdeaId)
		if err != nil {
			return
		}

		upvotes = append(upvotes, upvote)
	}

	return
}

func InsertUpvote(upvote types.Upvote) (err error) {
	_, err = Db.Exec("INSERT INTO upvotes(user_id, idea_id) VALUES($1, $2)", upvote.UserId, upvote.IdeaId)

	return
}

func CountUpvotes() (counts []types.UpvoteCount, err error) {
	rows, err := Db.Query(
		`SELECT ideas.id, count(upvotes.id) AS Count
			FROM ideas
			LEFT OUTER JOIN upvotes ON ideas.id = upvotes.idea_id
			GROUP BY ideas.id`)

	for rows.Next() {
		count := types.UpvoteCount{}

		err = rows.Scan(&count.IdeaId, &count.Count)
		if err != nil {
			return
		}

		counts = append(counts, count)
	}

	return
}

func UpdateUpvotesCount(counts []types.UpvoteCount, err error) {
	values := ""
	for _, count := range counts {
		values += fmt.Sprintf("(%v, %v),", count.IdeaId, count.Count)
	}

	values = strings.TrimRight(values, ",")

	q := `UPDATE ideas SET
	upvotes = v.count
	FROM (VALUES ` + values + `)
    AS v(id, count)
	WHERE ideas.id = v.id`

	_, err = Db.Exec(q)

	return
}

func InsertComment(comment types.Comment) (err error) {
	_, err = Db.Exec("INSERT INTO comments(idea_id, user_id, username, parent_id, comment,  created_at) VALUES($1, $2, $3, $4, $5, $6)", comment.IdeaId, comment.UserId, comment.Username, comment.ParentId, comment.Comment, comment.CreatedAt)

	return
}

func SelectComments(ideaId int64) (comments []*types.Comment, err error) {
	rows, err := Db.Query("SELECT id, idea_id, user_id, username, parent_id, comment, created_at FROM comments WHERE idea_id=$1", ideaId)
	if err != nil {
		return
	}

	for rows.Next() {
		comment := types.Comment{}

		err = rows.Scan(&comment.Id, &comment.IdeaId, &comment.UserId, &comment.Username, &comment.ParentId, &comment.Comment, &comment.CreatedAt)
		if err != nil {
			return
		}

		comments = append(comments, &comment)
	}

	return
}

func InsertVisit(visit types.Visit) (err error) {
	_, err = Db.Exec("INSERT INTO visits(path, visitor_id, city, country_code, continent_code, created_at) VALUES($1, $2, $3, $4, $5, $6)", visit.Path, visit.VisitorId, visit.City, visit.CountryCode, visit.ContinentCode, visit.CreatedAt)

	return
}

func GetVisitsCountSeriesByIdeaSlug(slug string) (visitsCountSeries types.VisitsCountSeries, err error) {
	rows, err := Db.Query(`SELECT count(id), to_char(date_trunc('day', (created_at)), 'MM-DD') as date
	FROM visits
	WHERE path=$1
	GROUP BY date
	ORDER BY date`, "/idea/" + slug)
	if err != nil {
		return
	}

	visitsCounts := types.VisitsCounts{}
	for rows.Next() {
		visitsCount := types.VisitsCount{}

		err = rows.Scan(&visitsCount.Count, &visitsCount.Date)
		if err != nil {
			return
		}

		visitsCounts = append(visitsCounts, visitsCount)
	}

	visitsCountSeries = visitsCounts.ToSeries()

	return
}

func GetVisitsCountSeriesByUsername(username string) (visitsCountSeries types.VisitsCountSeries, err error) {
	rows, err := Db.Query(`SELECT count(id), to_char(date_trunc('day', (created_at)), 'MM-DD') as date
	FROM visits
	WHERE path=$1
	GROUP BY date
	ORDER BY date`, "/user/" + username)
	if err != nil {
		return
	}

	visitsCounts := types.VisitsCounts{}
	for rows.Next() {
		visitsCount := types.VisitsCount{}

		err = rows.Scan(&visitsCount.Count, &visitsCount.Date)
		if err != nil {
			return
		}

		visitsCounts = append(visitsCounts, visitsCount)
	}

	visitsCountSeries = visitsCounts.ToSeries()

	return
}