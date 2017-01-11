package main

import (
	// "math/rand"
	"os"
	// "strconv"
	// "strings"
	// "time"

	"github.com/b4b4r07/slack-bot_test"
	// "github.com/b4b4r07/slack-bot_test/codename"
	"github.com/b4b4r07/slack-bot_test/gh"
	// "github.com/b4b4r07/slack-bot_test/healthcheck"
	// "github.com/b4b4r07/slack-bot_test/jokes"
	// "github.com/b4b4r07/slack-bot_test/quote"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
)

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	// create github client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	ghClient := github.NewClient(tc)

	// labelService := gh.NewGithubLabelService(ghClient)
	prService := gh.NewPullRequestService(ghClient)

	// create bot
	b := bot.Create(os.Getenv("SLACK_TOKEN"))

	/*
		b.Route("joke", func(args []string, ev *slack.MessageEvent) {
			b.SendRTM(jokes.GetRandomJoke().String(), ev.Channel)
		}, "will resond with a random joke")

		b.Route("codename", func(args []string, ev *slack.MessageEvent) {
			b.SendRTM("Here you go: "+codename.Get(), ev.Channel)
		}, "will generate a random codename")

		b.Route("motivate", func(args []string, ev *slack.MessageEvent) {
			b.SendRTM(quote.Get().Text, ev.Channel)
		}, "will motivate you!")

		b.Route("die", func(args []string, ev *slack.MessageEvent) {
			b.SendRTM(":(", ev.Channel)
			go b.Stop()
		}, "will shut me down")

		healthChecks := []*healthcheck.HealthCheck{}
		b.Route("healthcheck", func(args []string, ev *slack.MessageEvent) {
			if len(args) < 2 {
				b.SendRTM("Not enought arguments: <url> <period in ms>", ev.Channel)
				return
			}

			link := strings.Trim(args[0], "<>")

			for _, check := range healthChecks {
				if check.Target == link {
					b.SendRTM("Healthcheck for the link already exists", ev.Channel)
					return
				}
			}

			dur, err := strconv.Atoi(args[1])
			if err != nil {
				b.SendRTM("Bad duration number provided", ev.Channel)
				return
			}

			// dont check so fast
			if dur < 2000 {
				dur = 2000
			}

			c := healthcheck.NewHealthCheck(link, time.Duration(dur)*time.Millisecond)
			c.OnChange(func() {
				health := ""
				if c.Healthy {
					health = "Healthy"
				} else {
					health = "Down"
				}
				b.RTM().Reconnect()

				b.SendRTM("Target "+c.Target+" Health changed to "+health+"", ev.Channel)
			})

			b.SendRTM("Added health check for link: "+link, ev.Channel)

			c.Start()
			healthChecks = append(healthChecks, c)
		}, "will add healthcheck with period")

		b.Route("label", func(args []string, ev *slack.MessageEvent) {
			if len(args) < 2 {
				b.SendRTM("Two arguments must be specified: owner and repository", ev.Channel)
				return
			}

			owner := args[0]
			repo := args[1]

			labelService.RemoveAllLabels(owner, repo)
			labelService.CreateLabels(owner, repo, gh.DefaultLabels)

		}, "will add default labels to target repository")
	*/

	b.Route("pr2", func(args []string, ev *slack.MessageEvent) {
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
		build := "Help page incoming:"
		for _, route := range b.Router().Routes() {
			build += "\n`" + route.Trigger + "` => " + route.Description
		}

		b.SendRTM(build, ev.Channel)
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
