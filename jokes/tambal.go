package jokes

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

// TambalAPIURL is a jokes API URL
const TambalAPIURL = "http://tambal.azurewebsites.net/joke/random"

type TambalJoke struct {
	Joke string `json:"joke"`
}

func (j *TambalJoke) String() string {
	return j.Joke
}

type TambalJokeAPI struct{}

func (api *TambalJokeAPI) Get() Joke {
	req := gorequest.New()
	_, body, _ := req.Get(TambalAPIURL).End()

	joke := &TambalJoke{}
	json.Unmarshal([]byte(body), &joke)

	return joke
}
