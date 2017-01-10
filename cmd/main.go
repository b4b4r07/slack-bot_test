package main

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/b4b4r07/slack-bot_test"
	"github.com/b4b4r07/slack-bot_test/codename"
	"github.com/b4b4r07/slack-bot_test/gh"
	"github.com/b4b4r07/slack-bot_test/healthcheck"
	"github.com/b4b4r07/slack-bot_test/jokes"
	"github.com/b4b4r07/slack-bot_test/quote"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// create github client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	ghClient := github.NewClient(tc)

	labelService := gh.NewGithubLabelService(ghClient)

	// create bot
	bot := isaac.Create(os.Getenv("SLACK_TOKEN"))

	bot.Route("joke", func(args []string, ev *slack.MessageEvent) {
		bot.SendRTM(jokes.GetRandomJoke().String(), ev.Channel)
	}, "will resond with a random joke")

	bot.Route("codename", func(args []string, ev *slack.MessageEvent) {
		bot.SendRTM("Here you go: "+codename.Get(), ev.Channel)
	}, "will generate a random codename")

	bot.Route("motivate", func(args []string, ev *slack.MessageEvent) {
		bot.SendRTM(quote.Get().Text, ev.Channel)
	}, "will motivate you!")

	bot.Route("die", func(args []string, ev *slack.MessageEvent) {
		bot.SendRTM(":(", ev.Channel)
		go bot.Stop()
	}, "will shut me down")

	healthChecks := []*healthcheck.HealthCheck{}
	bot.Route("healthcheck", func(args []string, ev *slack.MessageEvent) {
		if len(args) < 2 {
			bot.SendRTM("Not enought arguments: <url> <period in ms>", ev.Channel)
			return
		}

		link := strings.Trim(args[0], "<>")

		for _, check := range healthChecks {
			if check.Target == link {
				bot.SendRTM("Healthcheck for the link already exists", ev.Channel)
				return
			}
		}

		dur, err := strconv.Atoi(args[1])
		if err != nil {
			bot.SendRTM("Bad duration number provided", ev.Channel)
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
			bot.RTM().Reconnect()

			bot.SendRTM("Target "+c.Target+" Health changed to "+health+"", ev.Channel)
		})

		bot.SendRTM("Added health check for link: "+link, ev.Channel)

		c.Start()
		healthChecks = append(healthChecks, c)
	}, "will add healthcheck with period")

	bot.Route("label", func(args []string, ev *slack.MessageEvent) {
		if len(args) < 2 {
			bot.SendRTM("Two arguments must be specified: owner and repository", ev.Channel)
			return
		}

		owner := args[0]
		repo := args[1]

		labelService.RemoveAllLabels(owner, repo)
		labelService.CreateLabels(owner, repo, gh.DefaultLabels)

	}, "will add default labels to target repository")

	bot.Route("help", func(args []string, ev *slack.MessageEvent) {
		build := "Help page incoming:"
		for _, route := range bot.Router().Routes() {
			build += "\n`" + route.Trigger + "` => " + route.Description
		}

		bot.SendRTM(build, ev.Channel)
	}, "will print help for all routes")

	err := bot.Run()
	if err != nil {
		panic(err)
	}
}
