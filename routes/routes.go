package routes

import (
	"github.com/21stio/go-ideahub/routes/templates"
	"github.com/21stio/go-ideahub/types"
	"github.com/21stio/go-ideahub/queries"
	"time"
	"fmt"
	"github.com/21stio/go-ideahub/utils"
	"sort"
)

const (
	HOME           = "/"
	AUTH           = "/auth"
	PUBLIC          = "/public"
	LOGIN          = "/login"
	LOGIN_CALLBACK = "/callback"
	LOGOUT         = "/logout"
	AUTHENTICATED  = "/u"
	SUBMIT         = "/submit"
	ME             = "/me"
	IDEA           = "/idea"
	USER           = "/user"
	UPVOTE         = "/upvote"
	COMMENT        = "/comment"
)

func addUserToBaseData(data templates.BaseData, user types.User) (baseData templates.BaseData) {
	baseData = data

	baseData.Title = user.Username
	baseData.LogoUrl = "http://" + user.AvatarUrl
	baseData.SubTitle = user.About
	baseData.TitleUrl = "/user/" + user.Username

	return baseData
}

func enrichIdeas(_ideas []types.Idea, baseData templates.BaseData) (ideas []types.Idea) {
	ideas = addAuthorsToIdeas(_ideas)

	ideas = addUpvotedToIdeas(ideas, baseData.Upvoted)

	ideas = addAgeLabelToIdeas(ideas)

	return
}

func addAuthorsToIdeas(_ideas []types.Idea) (ideas []types.Idea) {
	ideas = _ideas

	ids := []string{}

	for _, idea := range ideas {
		ids = append(ids, idea.CreatedBy)
	}

	ids = utils.GetUniqueStrings(ids)

	users, err := queries.SelectUsersByIds(ids)
	if err != nil {
		return
	}

	idToUsername := map[string]string{}

	for _, user := range users {
		idToUsername[user.Id] = user.Username
	}

	for i, idea := range ideas {
		ideas[i].Author = idToUsername[idea.CreatedBy]
	}

	return
}

func addUpvotedToIdeas(_ideas []types.Idea, upvoted map[int64]bool) (ideas []types.Idea) {
	ideas = _ideas

	for i, idea := range ideas {
		_, ok := upvoted[idea.Id]
		if ok {
			ideas[i].Upvoted = true
		}
	}

	return
}

func timeToLabel(t time.Time) (string) {
	since := time.Now().UTC().Sub(t)

	suffix := ""

	minutesSince := int64(since.Minutes())
	if minutesSince < 60 {
		if minutesSince == 1 {
			suffix = "minute"
		} else {
			suffix = "minutes"
		}

		return fmt.Sprintf("%v %v", minutesSince, suffix)
	}

	hoursSince := int64(since.Hours())
	if hoursSince < 24 {
		if hoursSince == 1 {
			suffix = "hour"
		} else {
			suffix = "hours"
		}

		return fmt.Sprintf("%v %v", hoursSince, suffix)
	} else if hoursSince < 48 {
		suffix = "day"

		return fmt.Sprintf("%v %v", 1, suffix)
	}

	return fmt.Sprintf("%v %v", hoursSince/24, suffix)
}

func addAgeLabelToIdeas(_ideas []types.Idea) (ideas []types.Idea) {
	ideas = _ideas

	for i, idea := range ideas {
		ideas[i].AgeLabel = timeToLabel(idea.CreatedAt)
	}

	return
}

func addAgeLabelToComments(_comments []*types.Comment) (comments []*types.Comment) {
	comments = _comments

	for i, comment := range comments {
		comments[i].AgeLabel = timeToLabel(comment.CreatedAt)
	}

	return
}

func addHierachyToComments(_comments []*types.Comment) (comments []*types.Comment) {
	comments = _comments

	idToComment := map[int64]*types.Comment{}

	for _, comment := range comments {
		idToComment[comment.Id] = comment
	}

	for _, comment := range comments {
		if comment.ParentId == 0 {
			continue
		}

		idToComment[comment.ParentId].Children = append(idToComment[comment.ParentId].Children, comment)
	}

	topLevelComments := []*types.Comment{}
	for _, comment := range comments {
		if comment.ParentId != 0 {
			continue
		}

		topLevelComments = append(topLevelComments, comment)
	}

	traverseComments(topLevelComments, 0)

	comments = topLevelComments

	return
}

func traverseComments(comments []*types.Comment, level int64) {
	sort.Sort(types.SortCommentsByCreatedAt(comments))

	for i, _ := range comments {
		comments[i].Level = level

		if len(comments[i].Children) != 0 {
			traverseComments(comments[i].Children, level+1)
		}
	}
}
