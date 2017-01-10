package isaac

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

type Isaac struct {
	slackRTM *slack.RTM
	slackAPI *slack.Client
	ghAPI    *github.Client

	router *Router

	run bool
}

func Create(slackKey string) *Isaac {
	slack := InitializeSlack(slackKey)
	gh := InitializeGithub()

	return &Isaac{
		slack.NewRTM(),
		slack,
		gh,
		NewRouter(slack),
		true,
	}
}

func (isaac *Isaac) RTM() *slack.RTM {
	return isaac.slackRTM
}

func (isaac *Isaac) SendRTM(msg, channel string) {
	isaac.RTM().SendMessage(isaac.RTM().NewOutgoingMessage(
		msg, channel,
	))
}

// Route adds a new route to the router of the bot
func (isaac *Isaac) Route(trigger string, call TriggerCall, description string) {
	isaac.router.Add(trigger, call, description)
}

func (isaac *Isaac) Router() *Router {
	return isaac.router
}

func (isaac *Isaac) Run() error {

	// start the connection manager
	go isaac.slackRTM.ManageConnection()

	for isaac.run {
		select {
		case msg := <-isaac.slackRTM.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				isaac.router.Match(ev)
			case *slack.ConnectedEvent:
				fmt.Println("Connected:", ev.Info.User.Name)
			default:
				// ignore
			}
		}
	}

	return nil
}

func (isaac *Isaac) Stop() {
	isaac.run = false
}
