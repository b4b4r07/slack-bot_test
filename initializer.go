package isaac

import (
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

// PanicIfErr panics the error if an error is passed
// to the function
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func InitializeSlack(key string) *slack.Client {
	api := slack.New(key)
	//api.SetDebug(true)

	return api
}

func InitializeGithub() *github.Client {
	return github.NewClient(nil)
}
