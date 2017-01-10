package jokes

import (
	"encoding/json"

	"github.com/parnurzeal/gorequest"
)

// ICNBAPIURL is jokes API URL
const ICNBAPIURL = "http://api.icndb.com/jokes/random"

type ICNDBJoke struct {
	Type  string `json:"type"`
	Value struct {
		ID         int      `json:"id"`
		Joke       string   `json:"joke"`
		Categories []string `json:"categories"`
	}
}

func (j *ICNDBJoke) String() string {
	return j.Value.Joke
}

type ICNBJokeAPI struct{}

func (api *ICNBJokeAPI) Get() Joke {
	req := gorequest.New()
	_, body, _ := req.Get(ICNBAPIURL).End()

	joke := &ICNDBJoke{}
	json.Unmarshal([]byte(body), &joke)

	return joke
}
