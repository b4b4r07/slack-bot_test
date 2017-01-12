package main

import (
	"fmt"
	"os"

	"github.com/b4b4r07/slack-bot_test"
	"github.com/b4b4r07/slack-bot_test/gh"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
)

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	ghClient := github.NewClient(tc)
	prService := gh.NewPullRequestService(ghClient)

	b := bot.Create(os.Getenv("SLACK_TOKEN"))
	b.Route("p-r", func(args []string, ev *slack.MessageEvent) {
		p := slack.PostMessageParameters{
			Username:  "pr-bot",
			IconEmoji: ":octocat:",
		}
		if len(args) < 1 {
			b.PostAttachment(ev.Channel, bot.Attachements(p, "too few argument", false))
			return
		}
		switch args[0] {
		case "list":
			p = prService.List("zplug", "zplug").MakeParams(p)
			if err := b.PostMessage(ev.Channel, "", p); err != nil {
				return
			}
			break
		default:
			b.PostAttachment(ev.Channel, bot.Attachements(p, "no such command", false))
			break
		}
	}, "pr-bot", bot.BotInfo{
		Name:  "pr-bot",
		Emoji: ":octocat:",
		Desc:  "github pr",
	})

	b.Route("help", func(args []string, ev *slack.MessageEvent) {
		params := slack.PostMessageParameters{}
		params.Attachments = []slack.Attachment{}
		params.Markdown = true
		for _, route := range b.Router().Routes() {
			params.Username = route.BotInfo.Name
			params.IconEmoji = route.BotInfo.Emoji
			params.Attachments = append(params.Attachments, slack.Attachment{
				Fallback:   "",
				AuthorName: fmt.Sprintf("%s %s", route.BotInfo.Emoji, route.BotInfo.Name),
				Title:      "",
				Text:       fmt.Sprintf("%s", route.Description),
				MarkdownIn: []string{"title", "text", "fields", "fallback"},
			})
		}
		b.PostAttachment(ev.Channel, params)
	}, "will print help for all routes", bot.BotInfo{
		Name:  "help-bot",
		Emoji: ":question:",
		Desc:  "help",
	})

	err := b.Run()
	if err != nil {
		panic(err)
	}
}
