package bot

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

type Bot struct {
	slackRTM *slack.RTM
	slackAPI *slack.Client
	ghAPI    *github.Client

	router *Router

	run bool
}

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

func (bot *Bot) SendRTM(msg, channel string) {
	bot.RTM().SendMessage(bot.RTM().NewOutgoingMessage(
		msg, channel,
	))
}

// Route adds a new route to the router of the bot
func (bot *Bot) Route(trigger string, call TriggerCall, description string) {
	bot.router.Add(trigger, call, description)
}

func (bot *Bot) Router() *Router {
	return bot.router
}

func (bot *Bot) Run() error {

	// start the connection manager
	go bot.slackRTM.ManageConnection()

	for bot.run {
		select {
		case msg := <-bot.slackRTM.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
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
