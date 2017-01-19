package bot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/b4b4r07/slack-bot_test/travis"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

type Bot struct {
	slackRTM *slack.RTM
	slackAPI *slack.Client
	ghAPI    *github.Client
	router   *Router
	run      bool
}

type Params *slack.PostMessageParameters

func Create(slackKey string) *Bot {
	slack := InitializeSlack(slackKey)
	gh := InitializeGithub()

	return &Bot{
		slack.NewRTM(),
		slack,
		gh,
		NewRouter(slack),
		true,
	}
}

func (bot *Bot) RTM() *slack.RTM {
	return bot.slackRTM
}

// TODO: irekae
func (bot *Bot) SendRTM(msg, channel string) {
	bot.RTM().SendMessage(bot.RTM().NewOutgoingMessage(
		msg, channel,
	))
}

func (bot *Bot) PostMessage(channel, msg string, params slack.PostMessageParameters) error {
	_, _, err := bot.slackAPI.PostMessage(channel, msg, params)
	return err
}

func (bot *Bot) PostAttachment(channel string, params slack.PostMessageParameters) error {
	_, _, err := bot.slackAPI.PostMessage(channel, "", params)
	return err
}

// func (bot *Bot) PostAttachment(channel, msg string, status bool) error {
// 	_, _, err := bot.slackAPI.PostMessage(channel, "", makeAttachements(msg, status))
// 	return err
// }

func Attachements(params slack.PostMessageParameters, msg string, status bool) slack.PostMessageParameters {
	color := "danger"
	if status {
		color = "good"
	}
	// params := slack.PostMessageParameters{
	// 	Markdown:  true,
	// 	Username:  "rc-bot",
	// 	IconEmoji: ":trollface:",
	// }
	params.Markdown = true
	params.Attachments = []slack.Attachment{}
	params.Attachments = append(params.Attachments, slack.Attachment{
		Fallback:   "",
		Title:      "",
		Text:       msg,
		MarkdownIn: []string{"title", "text", "fields", "fallback"},
		Color:      color,
	})
	return params
}

// Route adds a new route to the router of the bot
func (bot *Bot) Route(trigger string, call TriggerCall, description string, info BotInfo) {
	bot.router.Add(trigger, call, description, info)
}

func (bot *Bot) Router() *Router {
	return bot.router
}

func (bot *Bot) Hear(ev *slack.MessageEvent) error {
	var re *regexp.Regexp = regexp.MustCompile(`New comment by .*/pulls?/(\d+)`)
	for _, attachment := range ev.Attachments {
		if !strings.Contains(attachment.Text, "retest gotcha") {
			continue
		}
		pat := re.FindStringSubmatch(attachment.Pretext)
		if len(pat) < 2 {
			break
		}
		n, _ := strconv.Atoi(pat[1])
		if err := travis.RestartBuildFromPR("zplug/zplug", n); err != nil {
			return err
		}
	}
	return nil
}

func (bot *Bot) Run() error {
	// start the connection manager
	go bot.slackRTM.ManageConnection()

	for bot.run {
		select {
		case msg := <-bot.slackRTM.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				if ev.SubType == "bot_message" {
					bot.Hear(ev)
					break
				}
				bot.router.Match(ev)
			case *slack.ConnectedEvent:
				fmt.Println("Connected:", ev.Info.User.Name)
			default:
				// ignore
			}
		}
	}

	return nil
}

func (bot *Bot) Stop() {
	bot.run = false
}
