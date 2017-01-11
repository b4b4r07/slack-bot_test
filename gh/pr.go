package gh

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

var State = map[string]string{
	"open": "#67C63D",
}

// PullRequestService
type PullRequestService struct {
	client *github.Client
	issues []github.Issue
}

func NewPullRequestService(client *github.Client) *PullRequestService {
	return &PullRequestService{client, []github.Issue{}}
}

func (s *PullRequestService) List(owner, repo string) *PullRequestService {
	if owner == "" || repo == "" {
		return s
	}

	opt := &github.IssueListByRepoOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var issues []github.Issue
	for {
		repos, resp, err := s.client.Issues.ListByRepo(owner, repo, opt)
		if err != nil {
			return s
		}
		for _, v := range repos {
			issues = append(issues, *v)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return &PullRequestService{s.client, issues}
}

func (s *PullRequestService) MakeParams(p slack.PostMessageParameters) slack.PostMessageParameters {
	// var Params slack.PostMessageParameters = slack.PostMessageParameters{
	// 	Markdown:  true,
	// 	Username:  "pr-bot",
	// 	IconEmoji: ":octocat:",
	// }
	// p := Params
	p.Attachments = []slack.Attachment{}
	p.Markdown = true
	for _, issue := range s.issues {
		labels := []string{}
		if issue.PullRequestLinks == nil {
			continue
		}
		for _, label := range issue.Labels {
			labels = append(labels, "`"+*label.Name+"`")
		}
		p.Attachments = append(p.Attachments, slack.Attachment{
			Fallback:   fmt.Sprintf("%d - %s", *issue.Number, *issue.Title),
			Title:      fmt.Sprintf("<%s|#%d> %s", *issue.HTMLURL, *issue.Number, *issue.Title),
			Text:       strings.Join(labels, ", "),
			MarkdownIn: []string{"title", "text", "fields", "fallback"},
			Color:      State["open"],
			AuthorIcon: *issue.User.AvatarURL,
			AuthorName: "@" + *issue.User.Login,
			AuthorLink: *issue.User.HTMLURL,
		})
	}
	return p
}
